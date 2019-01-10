package router

import (
	"github.com/kernelgarden/diet/controllers"
	"github.com/labstack/echo"
)

func InitRoutes(e *echo.Echo) {
	controllers.HelloController{}.Init(e.Group("/hello"))
}