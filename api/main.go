package api

import (
	"html/template"
	"io"
	"os"

	"github.com/TimRobillard/todo_go/store"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type Api struct {
	store *store.PostgresStore
}

func Register(e *echo.Echo, pg *store.PostgresStore) error {
	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	if os.Getenv("ENV") == "development" {
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "method=${method}, uri=${uri}, status=${status}\n",
		}))
	} else {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.Recover())

	e.Renderer = t
	e.Static("/dist", "dist")

	RegisterTodoRoutes(e, pg)

	return nil
}
