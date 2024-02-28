package api

import (
	"net/http"
	"strconv"

	"github.com/TimRobillard/todo_go/api/middleware"
	"github.com/TimRobillard/todo_go/store"
	"github.com/TimRobillard/todo_go/views"
	"github.com/labstack/echo/v4"
)

type CustomContext struct {
	userId int
	echo.Context
}

func RegisterTodoRoutes(e *echo.Echo, pg store.TodoStorage) error {
	t := e.Group("/todo")
	t.Use(middleware.MyJwtMiddleware)

	t.POST("", func(c echo.Context) error {
		title := c.FormValue("todo")
		userId := middleware.GetUserIdFromRequest(c)
		todo, err := pg.CreateTodo(title, userId)

		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		component := views.Todo(todo)
		return render(c, component)
	})

	t.PUT("/:id", func(c echo.Context) error {
		_id := c.Param("id")
		id, err := strconv.Atoi(_id)
		userId := middleware.GetUserIdFromRequest(c)

		if err != nil {
			return err
		}

		err = pg.ToggleComplete(id, userId)
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		}

		return c.NoContent(http.StatusOK)
	})

	return nil
}
