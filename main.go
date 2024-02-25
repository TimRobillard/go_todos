package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"sort"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type Todo struct {
	Id       int
	Title    string
	Complete bool
}

var TODOS map[int]*Todo
var id int = 1

func main() {
	e := echo.New()

	TODOS = make(map[int]*Todo)

	TODOS[id] = &Todo{
		Id:       id,
		Title:    "My First Todo",
		Complete: false,
	}

	t := &Template{templates: template.Must(template.ParseGlob("views/*.html"))}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Renderer = t
	e.Static("/dist", "dist")

	e.GET("/", func(c echo.Context) error {
		var todos []*Todo
		var keys []int
		for k := range TODOS {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		for _, k := range keys {
			todo := TODOS[k]
			c := "false"
			if todo.Complete {
				c = "true"
			}

			fmt.Printf("i: %d, t: %s, c: %s\n", todo.Id, todo.Title, c)
			todos = append(todos, todo)
		}
		return c.Render(http.StatusOK, "index", todos)
	})

	e.POST("/todo", func(c echo.Context) error {
		id = id + 1
		title := c.FormValue("todo")
		todo := &Todo{
			Id:       id,
			Title:    title,
			Complete: false,
		}
		TODOS[id] = todo
		return c.Render(http.StatusOK, "todo", todo)
	})

	e.PUT("/todo/:id", func(c echo.Context) error {
		_id := c.Param("id")
		id, err := strconv.Atoi(_id)
		if err != nil {
			return err
		}
		co := TODOS[id].Complete
		t := "false"
		if co {
			t = "true"
		}
		fmt.Printf("complete? %s\n", t)
		TODOS[id].Complete = !TODOS[id].Complete

		return c.NoContent(http.StatusOK)
	})

	e.Logger.Fatal(e.Start(":4000"))
}
