package controllers

import (
	"github.com/kernelgarden/diet/constant"
	"github.com/kernelgarden/diet/factory"
	"github.com/kernelgarden/diet/models"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type FoodApiController struct {
}

func (f FoodApiController) Init(g *echo.Group) {
	g.GET("/test/:id", f.Test)
	g.GET("/:id", f.GetById)
	g.POST("", f.GetList)
	g.GET("/page", f.GetPage)
	g.PUT("/:id", f.Update)
	g.DELETE("", f.Delete)
}

func (FoodApiController) Test(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	var foods []models.BrandFood
	err = factory.DB().
		Join("INNER", "brand", "brand.id = food.brand_id").
		Where("brand.id = ?", id).
		Find(&foods)
	if err != nil {
		Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
	}

	return Success(ctx, foods)
}

func (FoodApiController) GetById(ctx echo.Context) error {
	param := ctx.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	food, err := models.Food{}.Get(id)
	if err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
	} else if food == nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.InExist))
	}

	return Success(ctx, food)
}

type FoodGetListInput struct {
	IdList	[]int64	`json:"id_list"`
}
type FoodGetListOutput struct {
	FoodList []*models.Food	`json:"food_list"`
}
func (FoodApiController) GetList(ctx echo.Context) error {
	var input FoodGetListInput
	if err := ctx.Bind(&input); err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	foodList := make([]*models.Food, len(input.IdList))
	for idx, id := range input.IdList {
		food, err := models.Food{}.Get(id)
		if err != nil {
			return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
		} else if food == nil {
			return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.InExist))
		}

		foodList[idx] = food
	}

	result := FoodGetListOutput{FoodList: foodList}

	return Success(ctx, result)
}

type FoodGetPageInput struct {
	Limit	int	`query:"limit"`
	Offset	int `query:"offset"`
}
type FoodGetPageOutput struct {
	FoodList []*models.Food `json:"food_list"`
}
func (FoodApiController) GetPage(ctx echo.Context) error {
	var input FoodGetPageInput
	if err := ctx.Bind(&input); err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	foodList, err := models.Food{}.GetAll(input.Offset, input.Limit)
	if err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
	}

	result := FoodGetListOutput{FoodList: foodList}

	return Success(ctx, result)
}

type FoodCreateInput struct {
}
type FoodCreateOutput struct {
	Food models.Food `json:"result"`
}
func (FoodApiController) Create(ctx echo.Context) error {

	return Success(ctx, nil)
}

type FoodDeleteInput struct {
	Id	int64	`query:"id"`
}
func (FoodApiController) Delete(ctx echo.Context) error {
	var input FoodDeleteInput
	if err := ctx.Bind(&input); err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	err := models.Food{}.Delete(input.Id)
	if err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
	}

	return Success(ctx, nil)
}

type FoodUpdateInput struct {
	CategoryId int64     `json:"category_id"`
	BrandId    int64     `json:"brand_id"`
	Name       string    `json:"name"`
	Weight     float64   `json:"weight"`
}
func (FoodApiController) Update(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	var input FoodUpdateInput
	if err := ctx.Bind(&input); err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	food, err := models.Food{}.Get(id)
	if err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
	} else if food == nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.InExist))
	}

	if input.CategoryId != 0 {
		food.CategoryId = input.CategoryId
	}

	if input.BrandId != 0 {
		food.BrandId = input.BrandId
	}

	if input.Name != "" {
		food.Name = input.Name
	}

	if input.Weight != 0 {
		food.Weight = input.Weight
	}

	if err = food.Update(); err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
	}

	return Success(ctx, nil)
}