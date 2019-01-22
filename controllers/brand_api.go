package controllers

import (
	"fmt"
	"github.com/kernelgarden/diet/constant"
	"github.com/kernelgarden/diet/factory"
	"github.com/kernelgarden/diet/models"
	"github.com/labstack/echo"
	"github.com/pangpanglabs/echoswagger"
	"net/http"
	"strconv"
)

type BrandApiController struct {
}

func (b BrandApiController) Init(g echoswagger.ApiGroup) {
	g.POST("", b.Create).
		AddParamQueryNested(BrandCreateInput{}).
		AddResponse(http.StatusOK, "생성된 brand의 정보를 반환합니다.", models.Brand{}, nil)
	g.GET("/dummy", b.CreateDummy)

	g.GET("/:id", b.GetById).
		AddParamQueryNested(BrandGetByIdInput{}).
		AddResponse(http.StatusOK, "조회할 brand의 정보를 반환합니다.", models.Brand{}, nil)
	g.POST("/list", b.GetList).
		AddParamBody(BrandGetListInput{}, "body", "", true).
		AddResponse(http.StatusOK, "조회할 brand 정보들의 리스트를 반환합니다.", BrandGetListOutput{}, nil)
	g.GET("/page", b.GetPage).
		AddParamQueryNested(BrandGetPageInput{}).
		AddResponse(http.StatusOK, "조회할 brand 정보들의 페이지를 반환합니다.", BrandGetPageOutput{}, nil)

	g.PUT("/:id", b.Update).
		AddParamQueryNested(BrandUpdateInput{}).
		AddResponse(http.StatusOK, "", nil, nil)

	g.DELETE("", b.Delete).
		AddParamQueryNested(BrandDeleteInput{}).
		AddResponse(http.StatusOK, "", nil, nil)
}

// TODO: add validator

type BrandGetByIdInput struct {
	Id int64 `query:"id" swagger:"desc(조회할 brand의 ID),required"`
}

func (BrandApiController) GetById(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
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
	IdList []int64 `json:"IdList" swagger:"desc(조회할 대상의 ID 리스트),required"`
}
type BrandGetListOutput struct {
	BrandList []*models.Brand `json:"brand_list"`
}

func (BrandApiController) GetList(ctx echo.Context) error {
	var input BrandGetListInput
	if err := ctx.Bind(&input); err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	brandList := make([]*models.Brand, len(input.IdList))
	for idx, id := range input.IdList {
		brand, err := models.Brand{}.Get(id)
		if err != nil {
			return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
		} else if brand == nil {
			return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.InExist))
		}

		brandList[idx] = brand
	}

	result := BrandGetListOutput{BrandList: brandList}

	return Success(ctx, result)
}

type BrandGetPageInput struct {
	Limit  int `query:"limit" swagger:"desc(조회할 brand의 개수),required"`
	Offset int `query:"offset" swagger:"desc(조회를 시작할 offset),required"`
}
type BrandGetPageOutput struct {
	BrandList []*models.Brand `json:"brand_list"`
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
	Name       string `json:"Name" swagger:"desc(등록할 이름),required"`
	ImgUrl     string `json:"ImgUrl" swagger:"desc(등록할 이미지 주소),required"`
	CategoryId int64  `json:"CategoryId" swagger:"desc(등록할 카테고리 ID),required"`
}

func (BrandApiController) Create(ctx echo.Context) error {
	var input BrandCreateInput
	if err := ctx.Bind(&input); err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	newBrand := models.Brand{Name: input.Name, ImgSrc: input.ImgUrl, CategoryId: input.CategoryId}
	_, err := newBrand.Create()
	if err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
	}

	return Success(ctx, newBrand)
}

type BrandDeleteInput struct {
	Id int64 `query:"id" swagger:"desc(삭제할 brand의 ID),required"`
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
	Name       string `json:"Name" swagger:"desc(변경할 이름(보내지 않으면 적용 X)),allowEmpty"`
	ImgSrc     string `json:"ImgUrl" swagger:"desc(변경할 이미지 주소(보내지 않으면 적용 X)),allowEmpty"`
	CategoryId int64  `json:"CategoryId" swagger:"desc(변경할 카테고리 ID(보내지 않으면 적용 X)),allowEmpty"`
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
