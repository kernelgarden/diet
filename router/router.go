package router

import (
	"github.com/kernelgarden/diet/controllers"
	"github.com/labstack/echo"
	"github.com/pangpanglabs/echoswagger"
)

func InitRoutes(e *echo.Echo) {
	controllers.HelloController{}.Init(e.Group("/hello"))

	r := echoswagger.New(e, "", "/doc", &echoswagger.Info{
		Title:	"Programmer-Diet",
		Description: "프로그래머는 다이어트를 어떻게 하는가에 대한 API 문서이다.",
		Version: "0.1",
	})

	controllers.BrandApiController{}.Init(r.Group("Brand", "/api/brands"))
	controllers.CategoryApiController{}.Init(r.Group("Category", "/api/categories"))
	controllers.FoodApiController{}.Init(r.Group("Food", "/api/foods"))
	controllers.NutrientApiController{}.Init(r.Group("Nutrient", "/api/nutrients"))
}
