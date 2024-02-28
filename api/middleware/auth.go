package middleware

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"github.com/TimRobillard/todo_go/util"
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
		cookie, err := c.Request().Cookie("_q")
		if err != nil {
			return c.Redirect(http.StatusTemporaryRedirect, "/login")
		}

		tokenStr := cookie.Value

		token, err := jwt.ParseWithClaims(
			tokenStr,
			&JwtCustomClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(util.GetEnv("JWT_SECRET", "secret")), nil
			},
		)
		if err != nil || !token.Valid {
			return c.Redirect(http.StatusTemporaryRedirect, "/login")
		}
		claims := token.Claims.(*JwtCustomClaims)
		c.Set("userId", claims.Id)
		return next(c)
	}
}

func GenerateToken(id int) (string, error) {
	claims := &JwtCustomClaims{
		id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(util.GetEnv("JWT_SECRET", "secret")))
	return t, err
}
