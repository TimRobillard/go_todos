package api

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"

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

	e.GET("/", func(c echo.Context) error {
		todos, err := pg.GetAllTodos()
		if err != nil {
			fmt.Printf(err.Error())
			return c.String(http.StatusNotFound, "Something went wrong")
		}

		return c.Render(http.StatusOK, "index", todos)
	})

	e.POST("/todo", func(c echo.Context) error {
		title := c.FormValue("todo")
		todo, err := pg.CreateToDo(title)

		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		}

		return c.Render(http.StatusOK, "todo", todo)
	})

	e.PUT("/todo/:id", func(c echo.Context) error {
		_id := c.Param("id")
		id, err := strconv.Atoi(_id)

		if err != nil {
			return err
		}

		err = pg.ToggleTodo(id)
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		}

		return c.NoContent(http.StatusOK)
	})

	return nil
}
