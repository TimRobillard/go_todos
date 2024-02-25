package middleware

import (
	"net/http"
	"strings"

	"github.com/TimRobillard/todo_go/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	Id int `json:"id"`
	jwt.RegisteredClaims
}

func GetUserIdFromRequest(c echo.Context) int {
	userId := c.Get("userId").(int)
	return userId
}

func MyJwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenStr := strings.Replace(c.Request().Header.Get("Authorization"), "Bearer ", "", 1)
		token, err := jwt.ParseWithClaims(tokenStr, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(util.GetEnv("JWT_SECRET", "secret")), nil
		})
		if err != nil || !token.Valid {
			return c.Redirect(http.StatusTemporaryRedirect, "/login")
		}
		claims := token.Claims.(*JwtCustomClaims)
		c.Set("userId", claims.Id)
		return next(c)
	}
}
