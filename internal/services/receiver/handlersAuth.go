package receiver

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"web-crawler/internal/models"
)

type inRegister struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *Service) Register(c echo.Context) error {
	var in inRegister

	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, errMsg{Message: err.Error()})
	}

	if in.Username == "" || in.Password == "" {
		return c.JSON(http.StatusBadRequest, errMsg{Message: "invalid json"})
	}

	zap.S().Debug("validating credentials")

	if !isUsernameCorrect(in.Username) {
		return c.JSON(http.StatusBadRequest, errMsg{Message: "login is not valid"})
	}

	if !isPasswordCorrect(in.Password) {
		return c.JSON(http.StatusBadRequest, errMsg{Message: "password is not valid"})
	}

	if !isEmailCorrect(in.Email) {
		return c.JSON(http.StatusBadRequest, errMsg{Message: "email is not valid"})
	}

	zap.S().Debug("checking if user exists")

	_, err := r.db.GetUserByUsername(in.Username)
	if err == nil {
		return c.JSON(http.StatusBadRequest, errMsg{Message: "user with this username already exists"})
	}

	_, err = r.db.GetUserByEmail(in.Email)
	if err == nil {
		return c.JSON(http.StatusBadRequest, errMsg{Message: "user with this email already exists"})
	}

	// creating access token
	tokenString, err := makeToken(in.Username, r.secretSignature)
	if err != nil {
		zap.S().Error("error while creating token: ", err)
		return echo.ErrInternalServerError
	}

	// creating hash to store password
	hash, err := generatePasswordHash(in.Password)
	if err != nil {
		zap.S().Error("error while hashing password: ", err)
		return echo.ErrInternalServerError
	}

	user := &models.User{
		Username: in.Username,
		Email:    in.Email,
		Password: hash,
	}

	_, err = r.db.AddUser(user)
	if err != nil {
		zap.S().Error("error while adding user: ", err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, struct {
		Access string `json:"access"`
	}{Access: tokenString})
}

type inLogin struct {
	Login    string `json:"login"` // logic can be email and username
	Password string `json:"password"`
}

func (r *Service) Login(c echo.Context) error {
	var in inLogin

	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, errMsg{Message: err.Error()})
	}

	var user *models.User

	user, err := r.db.GetUserByUsername(in.Login)
	if err != nil {
		zap.S().Debug("user not found by username")
	}

	if user == nil {
		user, err = r.db.GetUserByEmail(in.Login)
		if err != nil {
			zap.S().Debug("user not found by email")
		}
	}

	if user == nil {
		return c.JSON(http.StatusBadRequest, errMsg{Message: "user not found"})
	}

	err = comparePasswordWithHash(user.Password, in.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errMsg{Message: "invalid password"})
	}

	tokenString, err := makeToken(user.Username, r.secretSignature)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errMsg{Message: "error while creating token"})
	}

	return c.JSON(http.StatusOK, struct {
		Access string `json:"access"`
	}{Access: tokenString})
}

func (r *Service) GetUser(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		return c.JSON(http.StatusBadRequest, errMsg{Message: "invalid user"})
	}

	return c.JSON(http.StatusOK, struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}{user.Username, user.Email})
}

type inUpdateUser struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	NewPassword     string `json:"new_password"`
	CurrentPassword string `json:"current_password"`
}

func (r *Service) UpdateUser(c echo.Context) error {
	var in inUpdateUser

	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, errMsg{Message: err.Error()})
	}

	user, ok := c.Get("user").(*models.User)
	if !ok {
		return c.JSON(http.StatusBadRequest, errMsg{Message: "invalid user"})
	}

	if err := comparePasswordWithHash(user.Password, in.CurrentPassword); err != nil {
		return c.JSON(http.StatusBadRequest, errMsg{Message: "invalid password"})
	}

	if in.Username != "" {
		if !isUsernameCorrect(in.Username) {
			return c.JSON(http.StatusBadRequest, errMsg{Message: "login is not valid"})
		}
		user.Username = in.Username
	}

	if in.Email != "" {
		if !isEmailCorrect(in.Email) {
			return c.JSON(http.StatusBadRequest, errMsg{Message: "email is not valid"})
		}
		user.Email = in.Email
	}

	if in.NewPassword != "" {
		if !isPasswordCorrect(in.NewPassword) {
			return c.JSON(http.StatusBadRequest, errMsg{Message: "password is not valid"})
		}
		hash, err := generatePasswordHash(in.NewPassword)
		if err != nil {
			zap.S().Error("error while hashing password: ", err)
			return echo.ErrInternalServerError
		}
		user.Password = hash
	}

	if err := r.db.UpdateUser(user); err != nil {
		zap.S().Error("error while updating user: ", err)
		return echo.ErrInternalServerError
	}

	tokenString, err := makeToken(user.Username, r.secretSignature)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errMsg{Message: "error while creating token"})
	}

	return c.JSON(http.StatusOK, struct {
		Access string `json:"access"`
	}{Access: tokenString})
}
