package receiver

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
	"web-crawler/internal/config"
	"web-crawler/internal/models"
	"web-crawler/internal/services/receiver"
	"web-crawler/internal/utils"
)

var (
	cfg *config.Config
)

func getTestUserStruct() *models.User {
	return &models.User{
		ID: "00000000-0000-0000-0000-000000000000",
	}
}

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

	prj1 := struct {
		Name          string `json:"name"`
		StartUrl      string `json:"start_url"`
		NumberOfLinks int    `json:"number_of_links"`
		Depth         int    `json:"depth"`
	}{"newCreateProject" + strconv.Itoa(int(time.Now().Unix())), "https://google.com", 20, 20}

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/project/create",
		utils.GetReaderFromStruct(prj1),
	)
	req.Header.Set("Content-Type", "application/json; charset=utf8")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", getTestUserStruct())

	assert.NoError(t, r.CreateProject(c))

	assert.Equal(t, http.StatusOK, rec.Code)

	// wrong create data

	req = httptest.NewRequest(
		http.MethodPost,
		"/api/project/create",
		utils.GetReaderFromStruct(struct {
			Bib string `json:"bib"`
		}{"bib"}),
	)
	req.Header.Set("Content-Type", "application/json; charset=utf8")
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.Set("user", getTestUserStruct())

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
			Name          string `json:"name"`
			StartUrl      string `json:"start_url"`
			NumberOfLinks int    `json:"number_of_links"`
			Depth         int    `json:"depth"`
		}{"newGetProject" + strconv.Itoa(int(time.Now().Unix())), "https://google.com", 20, 20}),
	)
	createReq.Header.Set("Content-Type", "application/json; charset=utf8")
	createRec := httptest.NewRecorder()
	createCtx := e.NewContext(createReq, createRec)
	createCtx.Set("user", getTestUserStruct())

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
	getCtx.Set("user", getTestUserStruct())

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
	getCtx2.Set("user", getTestUserStruct())

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
			Name          string `json:"name"`
			StartUrl      string `json:"start_url"`
			NumberOfLinks int    `json:"number_of_links"`
			Depth         int    `json:"depth"`
		}{"newDeleteProject" + strconv.Itoa(int(time.Now().Unix())), "https://google.com", 20, 20}),
	)
	createReq.Header.Set("Content-Type", "application/json; charset=utf8")
	createRec := httptest.NewRecorder()
	createCtx := e.NewContext(createReq, createRec)
	createCtx.Set("user", getTestUserStruct())

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
	deleteCtx.Set("user", getTestUserStruct())

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
	prj1 := struct {
		Name          string `json:"name"`
		StartUrl      string `json:"start_url"`
		NumberOfLinks int    `json:"number_of_links"`
		Depth         int    `json:"depth"`
	}{"newGetAllProject" + strconv.Itoa(int(time.Now().Unix())), "https://google.com", 20, 20}

	createReq := httptest.NewRequest(
		http.MethodPost,
		"/project/create",
		utils.GetReaderFromStruct(prj1),
	)
	createReq.Header.Set("Content-Type", "application/json; charset=utf8")
	createRec := httptest.NewRecorder()
	createCtx := e.NewContext(createReq, createRec)
	createCtx.Set("user", getTestUserStruct())

	assert.NoError(t, r.CreateProject(createCtx))
	assert.Equal(t, http.StatusOK, createRec.Code)

	var outCreate struct {
		Id string `json:"id"`
	}
	assert.NoError(t, json.Unmarshal(createRec.Body.Bytes(), &outCreate))

	// Create 2
	prj2 := struct {
		Name          string `json:"name"`
		StartUrl      string `json:"start_url"`
		NumberOfLinks int    `json:"number_of_links"`
		Depth         int    `json:"depth"`
	}{"newGetAllProject2" + strconv.Itoa(int(time.Now().Unix())), "https://google.com", 20, 20}

	createReq2 := httptest.NewRequest(
		http.MethodPost,
		"/project/create",
		utils.GetReaderFromStruct(prj2),
	)
	createReq2.Header.Set("Content-Type", "application/json; charset=utf8")
	createRec2 := httptest.NewRecorder()
	createCtx2 := e.NewContext(createReq2, createRec2)
	createCtx2.Set("user", getTestUserStruct())

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
	getAllShortCtx.Set("user", getTestUserStruct())

	assert.NoError(t, r.GetAllShort(getAllShortCtx))
	assert.Equal(t, http.StatusOK, getAllShortRec.Code)

	assert.Contains(t, getAllShortRec.Body.String(), prj1.Name)
	assert.Contains(t, getAllShortRec.Body.String(), prj2.Name)
	assert.Contains(t, getAllShortRec.Body.String(), outCreate.Id)
	assert.Contains(t, getAllShortRec.Body.String(), outCreate2.Id)
}
