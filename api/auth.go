package api

import (
	"net/http"
	"time"

	"github.com/TimRobillard/todo_go/api/middleware"
	"github.com/TimRobillard/todo_go/store"
	"github.com/labstack/echo/v4"
)

func RegisterAuthRoutes(e *echo.Echo, pg store.UserStorage) error {
	e.GET("/register", func(c echo.Context) error {
		return c.Render(http.StatusOK, "register", nil)
	})
	e.GET("/login", func(c echo.Context) error {
		return c.Render(http.StatusOK, "login", nil)
	})

	g := e.Group("/auth")

	g.POST("/register", func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")
		if len(username) == 0 || len(password) == 0 {
			return c.String(http.StatusBadRequest, "username and password required")
		}

		user, err := pg.CreateUser(username, password)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		token, err := middleware.GenerateToken(user.Id)
		if err != nil {
			return c.String(http.StatusFound, "Incorrect username or password")
		}

		cookie := &http.Cookie{
			Value:   token,
			Name:    "_q",
			Path:    "/",
			Expires: time.Now().Add(24 * time.Hour),
		}
		c.SetCookie(cookie)
		c.Response().Header().Set("Hx-Redirect", "/")

		return c.String(http.StatusOK, "Success")
	})

	g.POST("/login", func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")
		if len(username) == 0 || len(password) == 0 {
			return c.String(http.StatusBadRequest, "username and password required")
		}

		user, err := pg.GetUserByUsername(username)
		if err != nil || !user.ValidatePassword(password) {
			return c.String(http.StatusBadRequest, err.Error())
		}

		token, err := middleware.GenerateToken(user.Id)
		if err != nil {
			return c.String(http.StatusFound, "Incorrect username or password")
		}

		cookie := &http.Cookie{
			Value:   token,
			Name:    "_q",
			Path:    "/",
			Expires: time.Now().Add(24 * time.Hour),
		}
		c.SetCookie(cookie)
		c.Response().Header().Set("Hx-Redirect", "/")

		return c.String(http.StatusOK, "Success")
	})

	return nil
}
