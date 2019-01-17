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
	g.POST("/page", b.GetPage)
	g.GET("/dummy", b.CreateDummy)
}

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
	Limit	int	`json:"limit"`
	Offset	int `json:"offset"`
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

func (BrandApiController) Create(ctx echo.Context) error {
	return Success(ctx, nil)
}

func (BrandApiController) Delete(ctx echo.Context) error {
	return Success(ctx, nil)
}

func (BrandApiController) Update(ctx echo.Context) error {
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