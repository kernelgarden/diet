package router

import (
	"github.com/kernelgarden/diet/controllers"
	"github.com/labstack/echo"
)

func InitRoutes(e *echo.Echo) {
	controllers.HelloController{}.Init(e.Group("/hello"))

	apiGroup := e.Group("/api")
	controllers.BrandApiController{}.Init(apiGroup.Group("/brands"))
	controllers.CategoryApiController{}.Init(apiGroup.Group("/categories"))
	controllers.FoodApiController{}.Init(apiGroup.Group("/foods"))
	controllers.NutrientApiController{}.Init(apiGroup.Group("/nutrients"))
}
