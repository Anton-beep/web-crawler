package receiver

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
	"web-crawler/internal/models"
	"web-crawler/internal/services/receiver"
	"web-crawler/internal/utils"
)

func TestRegisterUser(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	e := echo.New()
	r := receiver.New(1234, "../../../configs/.env")

	user1 := models.User{Username: "newUserRegister" + strconv.Itoa(int(time.Now().Unix())),
		Email:    "userRegister" + strconv.Itoa(int(time.Now().Unix())) + "@bib.com",
		Password: "Password123Password"}

	req1 := httptest.NewRequest(
		http.MethodPost,
		"/api/user/register",
		utils.GetReaderFromStruct(user1),
	)
	req1.Header.Set("Content-Type", "application/json; charset=utf8")
	rec := httptest.NewRecorder()
	c := e.NewContext(req1, rec)

	assert.NoError(t, r.Register(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	// login using username
	req2 := httptest.NewRequest(
		http.MethodPost,
		"/api/login",
		utils.GetReaderFromStruct(struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}{user1.Username, user1.Password}),
	)
	req2.Header.Set("Content-Type", "application/json; charset=utf8")
	rec = httptest.NewRecorder()
	c = e.NewContext(req2, rec)

	assert.NoError(t, r.Login(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	// login using email
	req3 := httptest.NewRequest(
		http.MethodPost,
		"/api/login",
		utils.GetReaderFromStruct(struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}{user1.Email, user1.Password}),
	)
	req3.Header.Set("Content-Type", "application/json; charset=utf8")
	rec = httptest.NewRecorder()
	c = e.NewContext(req3, rec)

	assert.NoError(t, r.Login(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	// invalid password
	req4 := httptest.NewRequest(
		http.MethodPost,
		"/api/register",
		utils.GetReaderFromStruct(struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}{Username: "newUserRegisterInvalid" + strconv.Itoa(int(time.Now().Unix())),
			Email:    "userRegister" + strconv.Itoa(int(time.Now().Unix())) + "@bib.com",
			Password: "invalidpass"}),
	)
	req4.Header.Set("Content-Type", "application/json; charset=utf8")
	rec = httptest.NewRecorder()

	c = e.NewContext(req4, rec)

	assert.NoError(t, r.Register(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestLoginUser(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	e := echo.New()
	r := receiver.New(1234, "../../../configs/.env")

	// register
	user1 := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{"newUserLogin" + strconv.Itoa(int(time.Now().Unix())),
		"userLogin" + strconv.Itoa(int(time.Now().Unix())) + "@bib.com",
		"Password123Password"}

	req1 := httptest.NewRequest(
		http.MethodPost,
		"/api/user/register",
		utils.GetReaderFromStruct(user1),
	)
	req1.Header.Set("Content-Type", "application/json; charset=utf8")

	rec := httptest.NewRecorder()
	c := e.NewContext(req1, rec)

	assert.NoError(t, r.Register(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	// login using username
	req2 := httptest.NewRequest(
		http.MethodPost,
		"/api/user/login",
		utils.GetReaderFromStruct(struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}{user1.Username, user1.Password}),
	)

	req2.Header.Set("Content-Type", "application/json; charset=utf8")
	rec = httptest.NewRecorder()
	c = e.NewContext(req2, rec)

	assert.NoError(t, r.Login(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	// login using email
	req3 := httptest.NewRequest(
		http.MethodPost,
		"/api/user/login",
		utils.GetReaderFromStruct(struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}{user1.Email, user1.Password}),
	)

	req3.Header.Set("Content-Type", "application/json; charset=utf8")
	rec = httptest.NewRecorder()
	c = e.NewContext(req3, rec)

	assert.NoError(t, r.Login(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	// login with wrong password

	req4 := httptest.NewRequest(
		http.MethodPost,
		"/api/user/login",
		utils.GetReaderFromStruct(struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}{user1.Username, "wrongPassword123"}),
	)

	req4.Header.Set("Content-Type", "application/json; charset=utf8")
	rec = httptest.NewRecorder()
	c = e.NewContext(req4, rec)

	assert.NoError(t, r.Login(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetUser(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	e := echo.New()
	r := receiver.New(1234, "../../../configs/.env")

	// register
	user1 := models.User{Username: "newUserGetUser" + strconv.Itoa(int(time.Now().Unix())),
		Email:    "userGetUser" + strconv.Itoa(int(time.Now().Unix())) + "@bib.com",
		Password: "Password123Password"}

	req1 := httptest.NewRequest(
		http.MethodPost,
		"/api/user/register",
		utils.GetReaderFromStruct(user1),
	)
	req1.Header.Set("Content-Type", "application/json; charset=utf8")

	rec := httptest.NewRecorder()
	c := e.NewContext(req1, rec)

	assert.NoError(t, r.Register(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	// login using username
	req2 := httptest.NewRequest(
		http.MethodPost,
		"/api/user/login",
		utils.GetReaderFromStruct(struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}{user1.Username, user1.Password}),
	)

	req2.Header.Set("Content-Type", "application/json; charset=utf8")
	rec = httptest.NewRecorder()
	c = e.NewContext(req2, rec)

	assert.NoError(t, r.Login(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	// get user
	req3 := httptest.NewRequest(
		http.MethodGet,
		"/api/user",
		nil,
	)

	rec = httptest.NewRecorder()
	c = e.NewContext(req3, rec)
	c.Set("user", &user1)

	assert.NoError(t, r.GetUser(c))
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthMiddleware(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	e := echo.New()
	r := receiver.New(1234, "../../../configs/.env")

	// register
	user1 := models.User{Username: "newUserAuthMiddleware" + strconv.Itoa(int(time.Now().Unix())),
		Email:    "userAuthMiddleware" + strconv.Itoa(int(time.Now().Unix())) + "@bib.com",
		Password: "Password123Password"}

	req1 := httptest.NewRequest(
		http.MethodPost,
		"/api/user/register",
		utils.GetReaderFromStruct(user1),
	)
	req1.Header.Set("Content-Type", "application/json; charset=utf8")

	rec := httptest.NewRecorder()
	c := e.NewContext(req1, rec)

	assert.NoError(t, r.Register(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	out := struct {
		Access string `json:"access"`
	}{}
	err := json.Unmarshal(rec.Body.Bytes(), &out)
	assert.NoError(t, err)

	// check middleware
	c.Request().Header.Set("Authorization", "Bearer "+out.Access)

	assert.NoError(t, r.AuthMiddleware(r.GetUser)(c))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, user1.Username, c.Get("user").(*models.User).Username)
}

func TestUpdateUser(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	e := echo.New()
	r := receiver.New(1234, "../../../configs/.env")

	// register
	user1 := models.User{
		Username: "newUserUpdate" + strconv.Itoa(int(time.Now().Unix())),
		Email:    "userUpdate" + strconv.Itoa(int(time.Now().Unix())) + "@bib.com",
		Password: "Password123Password",
	}

	req1 := httptest.NewRequest(
		http.MethodPost,
		"/api/user/register",
		utils.GetReaderFromStruct(user1),
	)
	req1.Header.Set("Content-Type", "application/json; charset=utf8")

	rec := httptest.NewRecorder()
	c := e.NewContext(req1, rec)

	assert.NoError(t, r.Register(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	access := struct {
		Access string `json:"access"`
	}{}
	err := json.Unmarshal(rec.Body.Bytes(), &access)
	assert.NoError(t, err)

	// update user
	updateData := struct {
		Username        string `json:"username"`
		Email           string `json:"email"`
		NewPassword     string `json:"new_password"`
		CurrentPassword string `json:"current_password"`
	}{
		Username:        "updatedUsername" + strconv.Itoa(int(time.Now().Unix())),
		Email:           "updatedEmail@bib.com",
		NewPassword:     "NewPassword123Password",
		CurrentPassword: user1.Password,
	}

	req3 := httptest.NewRequest(
		http.MethodPut,
		"/api/user",
		utils.GetReaderFromStruct(updateData),
	)
	req3.Header.Set("Content-Type", "application/json; charset=utf8")
	rec = httptest.NewRecorder()
	c = e.NewContext(req3, rec)

	c.Set("user", &user1)
	c.Request().Header.Set("Authorization", "Bearer "+access.Access)

	assert.NoError(t, r.AuthMiddleware(r.UpdateUser)(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	// verify updated user using login
	req4 := httptest.NewRequest(
		http.MethodPost,
		"/api/user/login",
		utils.GetReaderFromStruct(struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}{updateData.Username, updateData.NewPassword}),
	)
	req4.Header.Set("Content-Type", "application/json; charset=utf8")
	rec = httptest.NewRecorder()
	c = e.NewContext(req4, rec)

	assert.NoError(t, r.Login(c))
	assert.Equal(t, http.StatusOK, rec.Code)
}
