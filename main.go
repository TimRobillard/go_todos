package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/TimRobillard/todo_go/api"
	"github.com/TimRobillard/todo_go/store"
	"github.com/TimRobillard/todo_go/util"
)

func main() {
	e := echo.New()
	if os.Getenv("ENV") == "development" {
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "method=${method}, uri=${uri}, status=${status}\n",
		}))
	} else {
		e.Use(middleware.Logger())
	}

	pg, err := store.NewPostgresStore()
	if err != nil {
		e.Logger.Fatal(err.Error())
	}

	if err = pg.Init(); err != nil {
		e.Logger.Fatal(err.Error())
	}
	if err = pg.InitUser(); err != nil {
		e.Logger.Fatal(err.Error())
	}

	if err != nil {
		e.Logger.Fatal(err)
	}

	api.Register(e, pg)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", util.GetEnv("PORT", "4000"))))
}
