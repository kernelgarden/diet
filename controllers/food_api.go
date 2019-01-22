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
	g.POST("", f.Create).
		AddParamQueryNested(FoodCreateInput{}).
		AddResponse(http.StatusOK, "생성된 food의 정보를 반환합니다.", models.FoodJSON{}, nil)

	g.GET("/:id", f.GetById).
		AddParamQueryNested(FoodGetByIdInput{}).
		AddResponse(http.StatusOK, "조회할 food의 정보를 반환합니다.", models.FoodJSON{}, nil)
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
	Id int64 `query:"Id" swagger:"desc(조회할 food의 id),required"`
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

	result, err := food.ToJSON()
	if err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.InExist))
	}

	return Success(ctx, result)
}

type FoodGetListInput struct {
	IdList []int64 `json:"IdList" swagger:"desc(조회할 food의 ID 리스트),required"`
}
type FoodGetListOutput struct {
	FoodList []models.FoodJSON `json:"food_list"`
}

func (FoodApiController) GetList(ctx echo.Context) error {
	var input FoodGetListInput
	if err := ctx.Bind(&input); err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	foodList := make([]models.FoodJSON, len(input.IdList))
	for idx, id := range input.IdList {
		food, err := models.Food{}.Get(id)
		if err != nil {
			return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
		} else if food == nil {
			return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.InExist))
		}

		foodOutput, err := food.ToJSON()
		if err != nil {
			return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.InExist))
		}

		foodList[idx] = *foodOutput
	}

	result := FoodGetListOutput{FoodList: foodList}

	return Success(ctx, result)
}

type FoodGetPageInput struct {
	Limit  int `query:"limit" swagger:"desc(조회할 food의 개수),required"`
	Offset int `query:"offset" swagger:"desc(조회를 시작할 offset),required"`
}
type FoodGetPageOutput struct {
	FoodList []models.FoodJSON `json:"food_list"`
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

	foodGetByIdOutputList := make([]models.FoodJSON, len(foodList))
	for idx, food := range foodList {
		output, err := food.ToJSON()
		if err != nil {
			return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.Unknown))
		}

		foodGetByIdOutputList[idx] = *output
	}

	result := FoodGetPageOutput{FoodList: foodGetByIdOutputList}

	return Success(ctx, result)
}

type FoodCreateInput struct {
	CategoryId int64   `json:"CategoryId" swagger:"desc(생성할 food의 categoryId),required"`
	BrandId    int64   `json:"BrandId" swagger:"desc(생성할 food의 brandId),required"`
	Name       string  `json:"Name" swagger:"desc(생성할 food의 이름),required"`
	Weight     float64 `json:"Weight" swagger:"desc(생성할 food의 가중치),required"`

	Carbohydrate   float32 `json:"Carbohydrate" swagger:"desc(생성할 food의 탄수화물(g)),required"`
	Protein        float32 `json:"Protein" swagger:"desc(생성할 food의 단백질(g)),required"`
	SaturatedFat   float32 `json:"SaturatedFat" swagger:"desc(생성할 food의 포화지방(g)),required"`
	UnSaturatedFat float32 `json:"UnSaturatedFat" swagger:"desc(생성할 food의 불포화지방(g)),required"`
	TransFat       float32 `json:"TransFat" swagger:"desc(생성할 food의 트랜스지방(g)),required"`
	PerWeight      int32   `json:"PerWeight" swagger:"desc(생성할 food의 중량(g)),required"`
	Calorie        int64   `json:"Calorie" swagger:"desc(생성할 food의 칼로리(kcal)),required"`
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

	var category *models.Category
	var brand *models.Brand

	category, err = models.Category{}.Get(newFood.CategoryId)
	if err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
	} else if category == nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.InExist))
	}

	brand, err = models.Brand{}.Get(newFood.BrandId)
	if err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
	} else if category == nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.InExist))
	}

	result := models.FoodJSON{}.NewFoodJSON(newFood, newNutrient, *brand, *category)

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
	CategoryId int64   `json:"CategoryId" swagger:"desc(변경할 categoryId(보내지 않으면 적용X)),allowEmpty"`
	BrandId    int64   `json:"BrandId" swagger:"desc(변경할 brandId(보내지 않으면 적용X),allowEmpty"`
	Name       string  `json:"Name" swagger:"desc(변경할 이름(보내지 않으면 적용X),allowEmpty"`
	Weight     float64 `json:"Weight" swagger:"desc(변경할 가중치(보내지 않으면 적용X)),allowEmpty"`

	Carbohydrate   float32 `json:"Carbohydrate" swagger:"desc(생성할 food의 탄수화물(g)),allowEmpty"`
	Protein        float32 `json:"Protein" swagger:"desc(생성할 food의 단백질(g)),allowEmpty"`
	SaturatedFat   float32 `json:"SaturatedFat" swagger:"desc(생성할 food의 포화지방(g)),allowEmpty"`
	UnSaturatedFat float32 `json:"UnSaturatedFat" swagger:"desc(생성할 food의 불포화지방(g)),allowEmpty"`
	TransFat       float32 `json:"TransFat" swagger:"desc(생성할 food의 트랜스지방(g)),allowEmpty"`
	PerWeight      int32   `json:"PerWeight" swagger:"desc(생성할 food의 중량(g)),allowEmpty"`
	Calorie        int64   `json:"Calorie" swagger:"desc(생성할 food의 칼로리(kcal)),allowEmpty"`
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

	var nutrient *models.Nutrient
	nutrient, err = models.Nutrient{}.Get(id)
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
	/*
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
	*/
	return Success(ctx, nil)
}
