package controllers

import (
	"errors"
	"github.com/go-xorm/xorm"
	"github.com/kernelgarden/diet/constant"
	"github.com/kernelgarden/diet/factory"
	"github.com/kernelgarden/diet/models"
	"github.com/labstack/echo"
	"github.com/pangpanglabs/echoswagger"
	"net/http"
	"strconv"
)

type FoodApiController struct {
}

func (f FoodApiController) Init(g echoswagger.ApiGroup) {
	g.GET("/test/:id", f.Test)

	g.POST("", f.Create).
		AddParamQueryNested(FoodCreateInput{}).
		AddResponse(http.StatusOK, "생성된 food의 정보를 반환합니다.", FoodCreateOutput{}, nil)

	g.GET("/:id", f.GetById).
		AddParamQueryNested(FoodGetByIdInput{}).
		AddResponse(http.StatusOK, "조회할 food의 정보를 반환합니다.", FoodGetByIdOutput{}, nil)
	g.POST("/list", f.GetList).
		AddParamQueryNested(FoodGetListInput{}).
		AddResponse(http.StatusOK, "조회할 food 정보들의 리스트를 반환합니다.", FoodGetListOutput{}, nil)
	g.GET("/page", f.GetPage).
		AddParamQueryNested(FoodGetPageInput{}).
		AddResponse(http.StatusOK, "조회할 food 정보들의 페이지를 반환합니다.", FoodGetPageOutput{}, nil)

	g.PUT("/:id", f.Update).
		AddParamQueryNested(FoodUpdateInput{}).
		AddResponse(http.StatusOK, "", nil, nil)

	g.DELETE("", f.Delete).
		AddParamQueryNested(BrandDeleteInput{}).
		AddResponse(http.StatusOK, "", nil, nil)
}

type FoodGetByIdInput struct {
	Id int64 `json:"id" swagger:"desc(조회할 food의 id),required"`
}
type FoodGetByIdOutput struct {
	Food     models.Food     `json:"food"`
	Nutrient models.Nutrient `json:"nutrient"`
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

	var nutrient models.Nutrient
	if _, err := factory.DB().Where("food_id = ?", id).Get(&nutrient); err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.InExist))
	}

	result := FoodGetByIdOutput{Food: *food, Nutrient: nutrient}

	return Success(ctx, result)
}

type FoodGetListInput struct {
	IdList []int64 `json:"id_list" swagger:"desc(조회할 food의 ID 리스트),required"`
}
type FoodGetListOutput struct {
	FoodList []FoodGetByIdOutput `json:"food_list"`
}

func (FoodApiController) GetList(ctx echo.Context) error {
	var input FoodGetListInput
	if err := ctx.Bind(&input); err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	foodList := make([]FoodGetByIdOutput, len(input.IdList))
	for idx, id := range input.IdList {
		food, err := models.Food{}.Get(id)
		if err != nil {
			return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
		} else if food == nil {
			return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.InExist))
		}

		var nutrient models.Nutrient
		if has, err := factory.DB().Where("food_id = ?", food.Id).Get(&nutrient); err != nil {
			return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
		} else if !has {
			return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.InExist))
		}

		foodOutput := FoodGetByIdOutput{Food: *food, Nutrient: nutrient}

		foodList[idx] = foodOutput
	}

	result := FoodGetListOutput{FoodList: foodList}

	return Success(ctx, result)
}

type FoodGetPageInput struct {
	Limit  int `query:"limit" swagger:"desc(조회할 food의 개수),required"`
	Offset int `query:"offset" swagger:"desc(조회를 시작할 offset),required"`
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

	//TODO: change this
	result := FoodGetPageOutput{FoodList: foodList}

	return Success(ctx, result)
}

type FoodCreateInput struct {
	CategoryId int64   `json:"category_id" swagger:"desc(생성할 food의 categoryId),required"`
	BrandId    int64   `json:"brand_id" swagger:"desc(생성할 food의 brandId),required"`
	Name       string  `json:"name" swagger:"desc(생성할 food의 이름),required"`
	Weight     float64 `json:"weight" swagger:"desc(생성할 food의 가중치),required"`

	Carbohydrate   float32 `json:"carbohydrate" swagger:"desc(생성할 food의 탄수화물(g)),required"`
	Protein        float32 `json:"protein" swagger:"desc(생성할 food의 단백질(g)),required"`
	SaturatedFat   float32 `json:"saturated_fat" swagger:"desc(생성할 food의 포화지방(g)),required"`
	UnSaturatedFat float32 `json:"unsaturated_fat" swagger:"desc(생성할 food의 불포화지방(g)),required"`
	TransFat       float32 `json:"trans_fat" swagger:"desc(생성할 food의 트랜스지방(g)),required"`
	PerWeight      int32   `json:"per_weight" swagger:"desc(생성할 food의 중량(g)),requried"`
	Calorie        int64   `json:"calorie" swagger:"desc(생성할 food의 칼로리(kcal)),required"`
}
type FoodCreateOutput struct {
	Food     models.Food     `json:"food"`
	Nutrient models.Nutrient `json:"nutrient"`
}

func (FoodApiController) Create(ctx echo.Context) error {
	var input FoodCreateInput
	if err := ctx.Bind(&input); err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	newFood := models.Food{CategoryId: input.CategoryId, BrandId: input.BrandId, Name: input.Name, Weight: input.Weight}
	newNutrient := models.Nutrient{Carbohydrate: input.Carbohydrate, Protein: input.Protein, SaturatedFat: input.SaturatedFat,
		UnSaturatedFat: input.UnSaturatedFat, TransFat: input.TransFat, PerWeight: input.PerWeight, Calorie: input.Calorie}

	err := factory.Transaction(func(session *xorm.Session) error {
		if _, err := newFood.CreateWithSes(session); err != nil {
			return errors.New("food insert 실패")
		}

		newNutrient.FoodId = newFood.Id

		if _, err := newNutrient.CreateWithSes(session); err != nil {
			return errors.New("nutrient insert 실패")
		}

		return nil
	})
	if err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
	}

	result := FoodCreateOutput{Food: newFood, Nutrient: newNutrient}

	return Success(ctx, result)
}

type FoodDeleteInput struct {
	Id int64 `query:"id"`
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
	CategoryId int64   `json:"category_id" swagger:"desc(변경할 categoryId(보내지 않으면 적용X)),allowEmpty"`
	BrandId    int64   `json:"brand_id" swagger:"desc(변경할 brandId(보내지 않으면 적용X),allowEmpty"`
	Name       string  `json:"name" swagger:"desc(변경할 이름(보내지 않으면 적용X),allowEmpty"`
	Weight     float64 `json:"weight" swagger:"desc(변경할 가중치(보내지 않으면 적용X)),allowEmpty"`

	Carbohydrate   float32 `json:"carbohydrate" swagger:"desc(생성할 food의 탄수화물(g)),allowEmpty"`
	Protein        float32 `json:"protein" swagger:"desc(생성할 food의 단백질(g)),allowEmpty"`
	SaturatedFat   float32 `json:"saturated_fat" swagger:"desc(생성할 food의 포화지방(g)),allowEmpty"`
	UnSaturatedFat float32 `json:"unsaturated_fat" swagger:"desc(생성할 food의 불포화지방(g)),allowEmpty"`
	TransFat       float32 `json:"trans_fat" swagger:"desc(생성할 food의 트랜스지방(g)),allowEmpty"`
	PerWeight      int32   `json:"per_weight" swagger:"desc(생성할 food의 중량(g)),allowEmpty"`
	Calorie        int64   `json:"calorie" swagger:"desc(생성할 food의 칼로리(kcal)),allowEmpty"`
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

	var nutrient models.Nutrient
	_, err = factory.DB().Where("food_id = ?", id).Get(&nutrient)
	if err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.InExist))
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

	err = factory.Transaction(func(session *xorm.Session) error {
		if err = food.UpdateWithSes(session); err != nil {
			return errors.New("food update 실패")
		}

		if err = nutrient.UpdateWithSes(session); err != nil {
			return errors.New("nutrient update 실패")
		}

		return nil
	})
	if err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
	}

	return Success(ctx, nil)
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
