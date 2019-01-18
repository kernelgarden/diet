package controllers

import (
	"fmt"
	"github.com/kernelgarden/diet/constant"
	"github.com/kernelgarden/diet/factory"
	"github.com/kernelgarden/diet/models"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type BrandApiController struct {
}

func (b BrandApiController) Init(g *echo.Group) {
	g.GET("/:id", b.GetById)
	g.POST("", b.GetList)
	g.GET("/page", b.GetPage)
	g.GET("/dummy", b.CreateDummy)
	g.PUT("/:id", b.Update)
	g.DELETE("", b.Delete)
}

// TODO: add validator

func (BrandApiController) GetById(ctx echo.Context) error {
	param := ctx.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	brand, err := models.Brand{}.Get(id)
	if err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
	} else if brand == nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.InExist))
	}

	return Success(ctx, brand)
}

type BrandGetListInput struct {
	IdList	[]int64	`json:"id_list"`
}
type BrandGetListOutput struct {
	BrandList	[]*models.Brand `json:"brand_list"`
}
func (BrandApiController) GetList(ctx echo.Context) error {
	var input BrandGetListInput
	if err := ctx.Bind(&input); err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	brandList := make([]*models.Brand, 0)
	for _, id := range input.IdList {
		brand, err := models.Brand{}.Get(id)
		if err != nil {
			return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
		} else if brand == nil {
			return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.InExist))
		}

		brandList = append(brandList, brand)
	}

	result := BrandGetListOutput{BrandList: brandList}

	return Success(ctx, result)
}

type BrandGetPageInput struct {
	Limit	int	`query:"limit"`
	Offset	int `query:"offset"`
}
type BrandGetPageOutput struct {
	BrandList	[]*models.Brand	`json:"brand_list"`
}
func (BrandApiController) GetPage(ctx echo.Context) error {
	var input BrandGetPageInput
	if err := ctx.Bind(&input); err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	brandList, err := models.Brand{}.GetAll(input.Offset, input.Limit)
	if err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
	}

	result := BrandGetPageOutput{BrandList: brandList}

	return Success(ctx, result)
}

type BrandCreateInput struct {
}
type BrandCreateOutput struct {
	Brand 	models.Brand	`json:"result"`
}
func (BrandApiController) Create(ctx echo.Context) error {


	return Success(ctx, nil)
}

type BrandDeleteInput struct {
	Id	int64	`query:"id"`
}
func (BrandApiController) Delete(ctx echo.Context) error {
	var input BrandDeleteInput
	if err := ctx.Bind(&input); err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	err := models.Brand{}.Delete(input.Id)
	if err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
	}

	return Success(ctx, nil)
}

type BrandUpdateInput struct {
	Name       string    `json:"name"`
	ImgSrc     string    `json:"img_url"`
	CategoryId int64     `json:"category_id"`
}
func (BrandApiController) Update(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	var input BrandUpdateInput
	if err := ctx.Bind(&input); err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	brand, err := models.Brand{}.Get(id)
	if err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
	} else if brand == nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.InExist))
	}

	if input.Name != "" {
		brand.Name = input.Name
	}

	if input.ImgSrc != "" {
		brand.ImgSrc = input.ImgSrc
	}

	if input.CategoryId != 0 {
		brand.CategoryId = input.CategoryId
	}

	if err = brand.Update(); err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
	}

	return Success(ctx, nil)
}

func (BrandApiController) CreateDummy(ctx echo.Context) error {
	dummy := models.Brand{Name: "dummy", CategoryId: 0}

	_, err := dummy.Create()
	if err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
	}

	return Success(ctx, fmt.Sprintf("craete brand #%v", dummy.Id))
}