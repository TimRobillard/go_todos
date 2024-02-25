package api

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	myMiddleware "github.com/TimRobillard/todo_go/api/middleware"
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

	home := e.Group("/")
	home.Use(myMiddleware.MyJwtMiddleware)
	home.GET("/", func(c echo.Context) error {
		fmt.Println("index")
		todos, err := pg.GetAllTodos()
		if err != nil {
			fmt.Printf(err.Error())
			return c.String(http.StatusNotFound, "Something went wrong")
		}

		return c.Render(http.StatusOK, "index", todos)
	})

	RegisterAuthRoutes(e, pg)
	RegisterTodoRoutes(e, pg)

	return nil
}
