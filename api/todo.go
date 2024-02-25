package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TimRobillard/todo_go/store"
	"github.com/labstack/echo/v4"
)

func RegisterTodoRoutes(e *echo.Echo, pg *store.PostgresStore) error {
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
