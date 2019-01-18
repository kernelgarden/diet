package controllers

import (
	"github.com/kernelgarden/diet/constant"
	"github.com/kernelgarden/diet/factory"
	"github.com/kernelgarden/diet/models"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type CategoryApiController struct {
}

func (c CategoryApiController) Init(g *echo.Group) {
	g.GET("/:id", c.GetById)
	g.POST("", c.GetList)
	g.GET("/page", c.GetPage)
	g.PUT("/:id", c.Update)
	g.DELETE("", c.Delete)
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
	IdList	[]int64	`json:"id_list"`
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
	Limit	int	`query:"limit"`
	Offset	int `query:"offset"`
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
}
type CategoryCreateOutput struct {
	Brand 	models.Brand	`json:"result"`
}
func (CategoryApiController) Create(ctx echo.Context) error {

	return Success(ctx, nil)
}

type CategoryDeleteInput struct {
	Id	int64	`query:"id"`
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
	Name      string    `json:"name"`
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