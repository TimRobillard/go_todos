package api

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	myMiddleware "github.com/TimRobillard/todo_go/api/middleware"
	"github.com/TimRobillard/todo_go/store"
	"github.com/TimRobillard/todo_go/views"
)

type Api struct {
	store *store.PostgresStore
}

func Register(e *echo.Echo, pg *store.PostgresStore) error {
	e.Use(middleware.Recover())

	e.Static("/dist", "dist")

	home := e.Group("/")
	home.Use(myMiddleware.MyJwtMiddleware)

	home.GET("", func(c echo.Context) error {
		userId := myMiddleware.GetUserIdFromRequest(c)

		user, err := pg.GetUserById(userId)
		if err != nil {
			return c.String(http.StatusNotFound, "Something went wrong")
		}

		todos, err := pg.GetAllTodos(userId)
		if err != nil {
			return c.String(http.StatusNotFound, "Something went wrong")
		}

		component := views.Index(strings.Title(user.Username), todos)
		return render(c, component)
	})

	RegisterAuthRoutes(e, pg)
	RegisterTodoRoutes(e, pg)

	return nil
}
