package integrationTests

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"web-crauler/internal/services/receiver"
	"web-crauler/internal/utils"
)

func TestPing(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	assert.NoError(t, receiver.Pong(c))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "pong", rec.Body.String())
}

func TestCreateProject(t *testing.T) {
	e := echo.New()
	r := receiver.New(1234)

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
}

func TestGetProject(t *testing.T) {
	e := echo.New()
	r := receiver.New(1234)

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
		"/project/get",
		utils.GetReaderFromStruct(struct {
			Id string `json:"id"`
		}{outCreate.Id}),
	)
	getReq.Header.Set("Content-Type", "application/json; charset=utf8")
	getRec := httptest.NewRecorder()
	getCtx := e.NewContext(getReq, getRec)

	assert.NoError(t, r.GetProject(getCtx))
	assert.Equal(t, http.StatusOK, getRec.Code)
}

func TestDeleteProject(t *testing.T) {
	e := echo.New()
	r := receiver.New(1234)

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
		utils.GetReaderFromStruct(struct {
			Id string `json:"id"`
		}{outCreate.Id}),
	)
	deleteRec := httptest.NewRecorder()
	deleteCtx := e.NewContext(deleteReq, deleteRec)

	assert.NoError(t, r.DeleteProject(deleteCtx))
	assert.Equal(t, http.StatusOK, deleteRec.Code)

	// Verify project deletion
	getReq := httptest.NewRequest(
		http.MethodGet,
		"/project/get",
		utils.GetReaderFromStruct(struct {
			Id string `json:"id"`
		}{outCreate.Id}),
	)
	getRec := httptest.NewRecorder()
	getCtx := e.NewContext(getReq, getRec)

	assert.NoError(t, r.GetProject(getCtx))
	assert.Equal(t, http.StatusBadRequest, getRec.Code)
}
