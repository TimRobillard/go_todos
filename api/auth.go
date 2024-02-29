package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/TimRobillard/todo_go/api/middleware"
	"github.com/TimRobillard/todo_go/store"
	"github.com/TimRobillard/todo_go/views"
)

func RegisterAuthRoutes(e *echo.Echo, pg store.UserStorage) error {
	e.GET("/register", func(c echo.Context) error {
		component := views.Register()
		return render(c, component)
	})
	e.GET("/login", func(c echo.Context) error {
		component := views.Login()
		return render(c, component)
	})

	g := e.Group("/auth")

	g.POST("/register", func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")
		if len(username) == 0 || len(password) == 0 {
			return sendBadLogin(c, "Username and password required")
		}

		user, err := pg.CreateUser(username, password)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				c.Response().Status = 400
				component := views.UsernameTaken(username)
				return render(c, component)
			}
			return sendBadLogin(c, "Something went wrong, please try again")
		}

		token, err := middleware.GenerateToken(user.Id)
		if err != nil {
			return sendBadLogin(c, "Something went wrong, please try again")
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
			return sendBadLogin(c, "Username and password required")
		}

		user, err := pg.GetUserByUsername(username)
		if err != nil {
			return sendBadLogin(c, "Incorrect username or password")
		}

		valid := user.ValidatePassword(password)
		if valid == false {
			return sendBadLogin(c, "Incorrect username or password")
		}

		token, err := middleware.GenerateToken(user.Id)
		if err != nil {
			return sendBadLogin(c, "Incorrect username or password")
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

	g.POST("/logout", func(c echo.Context) error {
		cookie := &http.Cookie{
			Value:   "nil",
			Name:    "_q",
			Path:    "/",
			Expires: time.Now().Add(-1 * time.Minute),
		}
		c.SetCookie(cookie)
		c.Response().Header().Set("Hx-Redirect", "/login")
		return c.String(http.StatusTemporaryRedirect, "Success")
	})

	return nil
}

func sendBadLogin(c echo.Context, msg string) error {
	errorComponent := views.BadLogin(msg)
	c.Response().Status = 401
	return render(
		c, errorComponent,
	)
}
