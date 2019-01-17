package controllers

import (
	"github.com/kernelgarden/diet/constant"
	"github.com/labstack/echo"
	"net/http"
)

func Fail(ctx echo.Context, statusCode int, failResp constant.FailResp) error {
	return ctx.JSON(statusCode, failResp)
}

func Success(ctx echo.Context, resp interface{}) error {
	return ctx.JSON(http.StatusOK, resp)
}
