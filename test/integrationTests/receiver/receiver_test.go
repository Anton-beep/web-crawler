package receiver

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"web-crawler/internal/config"
	"web-crawler/internal/services/receiver"
	"web-crawler/internal/utils"
)

var (
	cfg *config.Config
)

func TestMain(m *testing.M) {
	cfg = config.NewConfig("../../../configs/.env")
	code := m.Run()
	os.Exit(code)
}

func TestCreateProject(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	e := echo.New()
	r := receiver.New(1234, "../../../configs/.env")

	req := httptest.NewRequest(
		http.MethodPost,
		"/project/create",
		utils.GetReaderFromStruct(struct {
			Name     string `json:"name"`
			StartUrl string `json:"start_url"`
		}{"newProject", "https://google.com"}),
	)
	req.Header.Set("Content-Type", "application/json; charset=utf8")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	assert.NoError(t, r.CreateProject(c))

	assert.Equal(t, http.StatusOK, rec.Code)

	// wrong create data

	req = httptest.NewRequest(
		http.MethodPost,
		"/project/create",
		utils.GetReaderFromStruct(struct {
			Bib string `json:"bib"`
		}{"bib"}),
	)
	req.Header.Set("Content-Type", "application/json; charset=utf8")
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	assert.NoError(t, r.CreateProject(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetProject(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	e := echo.New()
	r := receiver.New(1234, "../../../configs/.env")

	// Create
	createReq := httptest.NewRequest(
		http.MethodPost,
		"/project/create",
		utils.GetReaderFromStruct(struct {
			Name     string `json:"name"`
			StartUrl string `json:"start_url"`
		}{"newProject", "https://google.com"}),
	)
	createReq.Header.Set("Content-Type", "application/json; charset=utf8")
	createRec := httptest.NewRecorder()
	createCtx := e.NewContext(createReq, createRec)

	assert.NoError(t, r.CreateProject(createCtx))
	assert.Equal(t, http.StatusOK, createRec.Code)

	var outCreate struct {
		Id string `json:"id"`
	}
	assert.NoError(t, json.Unmarshal(createRec.Body.Bytes(), &outCreate))

	// Get project
	getReq := httptest.NewRequest(
		http.MethodGet,
		"/project/get/",
		nil,
	)
	getReq.Header.Set("Content-Type", "application/json; charset=utf8")
	getRec := httptest.NewRecorder()
	getCtx := e.NewContext(getReq, getRec)
	getCtx.SetParamNames("id")
	getCtx.SetParamValues(outCreate.Id)

	assert.NoError(t, r.GetProject(getCtx))
	assert.Equal(t, http.StatusOK, getRec.Code)
	assert.Contains(t, getRec.Body.String(), outCreate.Id)

	// Get project with wrong id
	getReq2 := httptest.NewRequest(
		http.MethodGet,
		"/project/get/",
		nil,
	)
	getReq2.Header.Set("Content-Type", "application/json; charset=utf8")
	getRec2 := httptest.NewRecorder()
	getCtx2 := e.NewContext(getReq2, getRec2)
	getCtx2.SetParamNames("id")
	getCtx2.SetParamValues("wrongId")

	assert.NoError(t, r.GetProject(getCtx2))
	assert.Equal(t, http.StatusBadRequest, getRec2.Code)
}

func TestDeleteProject(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	e := echo.New()
	r := receiver.New(1234, "../../../configs/.env")

	// Create project
	createReq := httptest.NewRequest(
		http.MethodPost,
		"/project/create",
		utils.GetReaderFromStruct(struct {
			Name     string `json:"name"`
			StartUrl string `json:"start_url"`
		}{"newProject", "https://google.com"}),
	)
	createReq.Header.Set("Content-Type", "application/json; charset=utf8")
	createRec := httptest.NewRecorder()
	createCtx := e.NewContext(createReq, createRec)

	assert.NoError(t, r.CreateProject(createCtx))
	assert.Equal(t, http.StatusOK, createRec.Code)

	var outCreate struct {
		Id string `json:"id"`
	}
	assert.NoError(t, json.Unmarshal(createRec.Body.Bytes(), &outCreate))

	// Delete project
	deleteReq := httptest.NewRequest(
		http.MethodDelete,
		"/project/delete",
		nil,
	)
	deleteRec := httptest.NewRecorder()
	deleteCtx := e.NewContext(deleteReq, deleteRec)
	deleteCtx.SetParamNames("id")
	deleteCtx.SetParamValues(outCreate.Id)

	assert.NoError(t, r.DeleteProject(deleteCtx))
	assert.Equal(t, http.StatusOK, deleteRec.Code)

	// Verify project deletion
	getReq := httptest.NewRequest(
		http.MethodGet,
		"/project/get",
		nil,
	)
	getRec := httptest.NewRecorder()
	getCtx := e.NewContext(getReq, getRec)
	getCtx.SetParamNames("id")
	getCtx.SetParamValues(outCreate.Id)

	assert.NoError(t, r.GetProject(getCtx))
	assert.Equal(t, http.StatusNotFound, getRec.Code)
}

func TestGetAllShort(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	e := echo.New()
	r := receiver.New(1234, "../../../configs/.env")

	// Create
	createReq := httptest.NewRequest(
		http.MethodPost,
		"/project/create",
		utils.GetReaderFromStruct(struct {
			Name     string `json:"name"`
			StartUrl string `json:"start_url"`
		}{"newProject", "https://google.com"}),
	)
	createReq.Header.Set("Content-Type", "application/json; charset=utf8")
	createRec := httptest.NewRecorder()
	createCtx := e.NewContext(createReq, createRec)

	assert.NoError(t, r.CreateProject(createCtx))
	assert.Equal(t, http.StatusOK, createRec.Code)

	var outCreate struct {
		Id string `json:"id"`
	}
	assert.NoError(t, json.Unmarshal(createRec.Body.Bytes(), &outCreate))

	// Create 2
	createReq2 := httptest.NewRequest(
		http.MethodPost,
		"/project/create",
		utils.GetReaderFromStruct(struct {
			Name     string `json:"name"`
			StartUrl string `json:"start_url"`
		}{"newProject2", "https://google.com"}),
	)
	createReq2.Header.Set("Content-Type", "application/json; charset=utf8")
	createRec2 := httptest.NewRecorder()
	createCtx2 := e.NewContext(createReq2, createRec2)

	assert.NoError(t, r.CreateProject(createCtx2))
	assert.Equal(t, http.StatusOK, createRec2.Code)

	var outCreate2 struct {
		Id string `json:"id"`
	}
	assert.NoError(t, json.Unmarshal(createRec2.Body.Bytes(), &outCreate2))

	// Get all short
	getAllShortReq := httptest.NewRequest(
		http.MethodGet,
		"/project/getAllShort",
		nil,
	)
	getAllShortRec := httptest.NewRecorder()
	getAllShortCtx := e.NewContext(getAllShortReq, getAllShortRec)

	assert.NoError(t, r.GetAllShort(getAllShortCtx))
	assert.Equal(t, http.StatusOK, getAllShortRec.Code)

	assert.Contains(t, getAllShortRec.Body.String(), "newProject")
	assert.Contains(t, getAllShortRec.Body.String(), "newProject2")
	assert.Contains(t, getAllShortRec.Body.String(), outCreate.Id)
	assert.Contains(t, getAllShortRec.Body.String(), outCreate2.Id)
}
