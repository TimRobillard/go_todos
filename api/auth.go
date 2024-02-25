package api

import (
	"net/http"

	"github.com/TimRobillard/todo_go/store"
	"github.com/labstack/echo/v4"
)

func RegisterAuthRoutes(e *echo.Echo, pg *store.PostgresStore) error {
	e.GET("/register", func(c echo.Context) error {
		return c.Render(http.StatusOK, "register", nil)
	})
	e.GET("/login", func(c echo.Context) error {
		return c.Render(http.StatusOK, "register", nil)
	})

	g := e.Group("/auth")

	g.POST("/register", func(c echo.Context) error {
		return c.String(http.StatusUnauthorized, "Incorrect username or password")
	})

	return nil
}
