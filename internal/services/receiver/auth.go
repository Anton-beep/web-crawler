package receiver

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

func makeToken(login string, secretSignature []byte) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": login,
		"nbf":  now.Unix(),
		"exp":  now.Add(5 * time.Minute).Unix(),
		"iat":  now.Unix(),
	})

	return token.SignedString(secretSignature)
}

func generatePasswordHash(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash, nil
}

func comparePasswordWithHash(existing string, incoming string) error {
	return bcrypt.CompareHashAndPassword([]byte(existing), []byte(incoming))
}

func (r *Service) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		access := c.Request().Header.Get("Authorization")
		access = strings.Replace(access, "Bearer ", "", 1)

		tokenFromString, err := jwt.Parse(access, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				panic(fmt.Errorf("unexpected signing method: %v", token.Header["alg"]))
			}

			return r.secretSignature, nil
		})

		if err != nil {
			zap.S().Debug("error while parsing token: ", err)
			return c.JSON(http.StatusUnauthorized, errMsg{Message: "invalid token"})
		}

		claims, ok := tokenFromString.Claims.(jwt.MapClaims)

		if !ok {
			zap.S().Debug("error while parsing claims")
			return c.JSON(http.StatusUnauthorized, errMsg{Message: "invalid token"})
		}

		user, err := r.db.GetUserByUsername(claims["name"].(string))
		if err != nil {
			zap.S().Debug("error while getting user: ", err)
			return c.JSON(http.StatusUnauthorized, errMsg{Message: "invalid token"})
		}

		c.Set("user", user)
		return next(c)
	}
}
