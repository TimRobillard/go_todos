package main

import (
	"fmt"

	"github.com/TimRobillard/todo_go/api"
	"github.com/TimRobillard/todo_go/store"
	"github.com/TimRobillard/todo_go/util"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	pg, err := store.NewPostgresStore()
	pg.Init()

	if err != nil {
		e.Logger.Fatal(err)
	}

	api.Register(e, pg)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", util.GetEnv("PORT", "4000"))))
}
