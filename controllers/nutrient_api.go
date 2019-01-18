package controllers

import (
	"github.com/kernelgarden/diet/constant"
	"github.com/kernelgarden/diet/factory"
	"github.com/kernelgarden/diet/models"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type NutrientApiController struct {
}

func (n NutrientApiController) Init(g *echo.Group) {
	g.GET("/:id", n.GetById)
	g.POST("", n.GetList)
	g.GET("/page", n.GetPage)
	g.PUT("/:id", n.Update)
	g.DELETE("", n.Delete)
}

func (NutrientApiController) GetById(ctx echo.Context) error {
	param := ctx.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	nutrient, err := models.Nutrient{}.Get(id)
	if err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
	} else if nutrient == nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.InExist))
	}

	return Success(ctx, nutrient)
}

type NutrientGetListInput struct {
	IdList	[]int64	`json:"id_list"`
}
type NutrientGetListOutput struct {
	NutrientList []*models.Nutrient `json:"nutrient_list"`
}
func (NutrientApiController) GetList(ctx echo.Context) error {
	var input NutrientGetListInput
	if err := ctx.Bind(&input); err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	nutrientList := make([]*models.Nutrient, len(input.IdList))
	for idx, id := range input.IdList {
		nutrient, err := models.Nutrient{}.Get(id)
		if err != nil {
			return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
		} else if nutrient == nil {
			return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.InExist))
		}

		nutrientList[idx] = nutrient
	}

	result := NutrientGetListOutput{NutrientList: nutrientList}

	return Success(ctx, result)
}

type NutrientGetPageInput struct {
	Limit	int	`query:"limit"`
	Offset	int `query:"offset"`
}
type NutrientGetPageOutput struct {
	NutrientList []*models.Nutrient `json:"nutrient_list"`
}
func (NutrientApiController) GetPage(ctx echo.Context) error {
	var input NutrientGetPageInput
	if err := ctx.Bind(&input); err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	nutrientList, err := models.Nutrient{}.GetAll(input.Offset, input.Limit)
	if err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
	}

	result := NutrientGetPageOutput{NutrientList: nutrientList}

	return Success(ctx, result)
}

type NutrientCreateInput struct {
}
type NutrientCreateOutput struct {
	Nutrient models.Nutrient `json:"result"`
}
func (NutrientApiController) Create(ctx echo.Context) error {

	return Success(ctx, nil)
}

type NutrientDeleteInput struct {
	Id	int64	`query:"id"`
}
func (NutrientApiController) Delete(ctx echo.Context) error {
	var input NutrientDeleteInput
	if err := ctx.Bind(&input); err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	err := models.Nutrient{}.Delete(input.Id)
	if err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
	}

	return Success(ctx, nil)
}

type NutrientUpdateInput struct {
	FoodId         int64     `json:"food_id"`
	Carbohydrate   float32   `json:"carbohydrate"`
	Protein        float32   `json:"protein"`
	SaturatedFat   float32   `json:"saturated_fat"`
	UnSaturatedFat float32   `json:"unsaturated_fat"`
	TransFat       float32   `json:"trans_fat"`
	PerWeight      int32     `json:"per_weight"`
	Calorie        int64     `json:"calorie"`
}
func (NutrientApiController) Update(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	var input NutrientUpdateInput
	if err := ctx.Bind(&input); err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	nutrient, err := models.Nutrient{}.Get(id)
	if err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
	} else if nutrient == nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.InExist))
	}

	if input.FoodId != 0 {
		nutrient.FoodId = input.FoodId
	}

	if input.Carbohydrate > 0 {
		nutrient.Carbohydrate = input.Carbohydrate
	}

	if input.Protein > 0 {
		nutrient.Protein = input.Protein
	}

	if input.SaturatedFat > 0 {
		nutrient.SaturatedFat = input.SaturatedFat
	}

	if input.UnSaturatedFat > 0 {
		nutrient.UnSaturatedFat = input.UnSaturatedFat
	}

	if input.TransFat > 0 {
		nutrient.TransFat = input.TransFat
	}

	if input.PerWeight != 0 {
		nutrient.PerWeight = input.PerWeight
	}

	if input.Calorie != 0 {
		nutrient.Calorie = input.Calorie
	}

	if err = nutrient.Update(); err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
	}

	return Success(ctx, nil)
}