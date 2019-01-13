package controllers

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

type HelloController struct {
	name string
}

func (c HelloController) Init(g *echo.Group) {
	g.GET("", c.Get)
	g.GET("/:name", c.GetByName)
}

func (HelloController) Get(ctx echo.Context) error {
	ctx.String(http.StatusOK, "Hello, anonymous")
	return nil
}

func (HelloController) GetByName(ctx echo.Context) error {
	name := ctx.Param("name")
	ctx.String(http.StatusOK, fmt.Sprintf("Hello, %s", name))
	return nil
}
