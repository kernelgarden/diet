package controllers

import (
	"github.com/kernelgarden/diet/constant"
	"github.com/kernelgarden/diet/factory"
	"github.com/kernelgarden/diet/models"
	"github.com/labstack/echo"
	"github.com/pangpanglabs/echoswagger"
	"net/http"
	"strconv"
)

type CategoryApiController struct {
}

func (c CategoryApiController) Init(g echoswagger.ApiGroup) {
	g.POST("", c.Create).
		AddParamQueryNested(CategoryCreateInput{}).
		AddResponse(http.StatusOK, "생성된 category의 정보를 반환합니다.", models.Category{}, nil)

	g.GET("/:id", c.GetById).
		AddParamQueryNested(CategoryGetByIdInput{}).
		AddResponse(http.StatusOK, "조회할 category의 정보를 반환합니다.", models.Category{}, nil)
	g.POST("/list", c.GetList).
		AddParamBody(CategoryGetListInput{}, "body", "", true).
		AddResponse(http.StatusOK, "조회할 category 정보들의 리스트를 반환합니다.", CategoryGetListOutput{}, nil)
	g.GET("/page", c.GetPage).
		AddParamQueryNested(CategoryGetPageInput{}).
		AddResponse(http.StatusOK, "조회할 category 정보들의 페이지를 반환합니다.", CategoryGetPageOutput{}, nil)

	g.PUT("/:id", c.Update).
		AddParamQueryNested(CategoryUpdateInput{}).
		AddResponse(http.StatusOK, "", nil, nil)

	g.DELETE("", c.Delete).
		AddParamQueryNested(CategoryDeleteInput{}).
		AddResponse(http.StatusOK, "", nil, nil)
}

type CategoryGetByIdInput struct {
	Id	int64	`query:"id" swagger:"desc(조회할 category의 ID),required"`
}
func (CategoryApiController) GetById(ctx echo.Context) error {
	param := ctx.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	category, err := models.Category{}.Get(id)
	if err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
	} else if category == nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.InExist))
	}

	return Success(ctx, category)
}

type CategoryGetListInput struct {
	IdList	[]int64	`json:"id_list" swagger:"desc(조회할 category ID 리스트),required"`
}
type CategoryGetListOutput struct {
	CategoryList	[]*models.Category `json:"category_list"`
}
func (CategoryApiController) GetList(ctx echo.Context) error {
	var input CategoryGetListInput
	if err := ctx.Bind(&input); err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	categoryList := make([]*models.Category, len(input.IdList))
	for idx, id := range input.IdList {
		category, err := models.Category{}.Get(id)
		if err != nil {
			return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
		} else if category == nil {
			return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.InExist))
		}

		categoryList[idx] = category
	}

	result := CategoryGetListOutput{CategoryList: categoryList}

	return Success(ctx, result)
}

type CategoryGetPageInput struct {
	Limit	int	`query:"limit" swagger:"desc(조회할 category의 개수),required"`
	Offset	int `query:"offset" swagger:"desc(조회를 시작할 offset),required"`
}
type CategoryGetPageOutput struct {
	CategoryList []*models.Category `json:"category_list"`
}
func (CategoryApiController) GetPage(ctx echo.Context) error {
	var input CategoryGetPageInput
	if err := ctx.Bind(&input); err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	categoryList, err := models.Category{}.GetAll(input.Offset, input.Limit)
	if err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
	}

	result := CategoryGetPageOutput{CategoryList: categoryList}

	return Success(ctx, result)
}

type CategoryCreateInput struct {
	Name      string    `json:"name" swagger:"desc(생성할 category의 이름),required"`
}
func (CategoryApiController) Create(ctx echo.Context) error {
	var input CategoryCreateInput
	if err := ctx.Bind(&input); err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	newCategory := models.Category{Name: input.Name}
	_, err := newCategory.Create()
	if err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
	}

	return Success(ctx, newCategory)
}

type CategoryDeleteInput struct {
	Id	int64	`query:"id" swagger:"desc(삭제할 category의 ID),required"`
}
func (CategoryApiController) Delete(ctx echo.Context) error {
	var input CategoryDeleteInput
	if err := ctx.Bind(&input); err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	err := models.Category{}.Delete(input.Id)
	if err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
	}

	return Success(ctx, nil)
}

type CategoryUpdateInput struct {
	Name      string    `json:"name" swagger:"desc(변경할 category 이름(보내지 않으면 적용X)),allowEmpty"`
}
func (CategoryApiController) Update(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	var input CategoryUpdateInput
	if err := ctx.Bind(&input); err != nil {
		return Fail(ctx, http.StatusBadRequest, factory.NewFailResp(constant.InvalidRequestFormat))
	}

	category, err := models.Category{}.Get(id)
	if err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
	} else if category == nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.InExist))
	}

	if input.Name != "" {
		category.Name = input.Name
	}

	if err = category.Update(); err != nil {
		return Fail(ctx, http.StatusInternalServerError, factory.NewFailResp(constant.Unknown))
	}

	return Success(ctx, nil)
}