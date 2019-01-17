package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/kernelgarden/diet/constant"
	"github.com/kernelgarden/diet/models"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kernelgarden/diet/factory"
	"github.com/kernelgarden/diet/router"
	"github.com/kernelgarden/goutils/config"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	// Read config file
	curPath, err := getCurPath()
	if err != nil {
		panic(err)
	}

	var c constant.Config
	err = config.Read(curPath, "config", &c)
	if err != nil {
		panic(err)
	}
	constant.GlobalCtx = context.WithValue(constant.GlobalCtx, constant.CtxConfig, &c)

	db := factory.DB()
	if db == nil {
		panic(err)
	}
	defer db.(*xorm.Engine).Close()
	Sync()

	e := echo.New()

	router.InitRoutes(e)

	e.Static("/static", "static")
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())

	e.Debug = c.Debug

	var port string
	if c.Httpport == "" {
		port = "3030"
	} else {
		port = c.Httpport
	}

	if err := e.Start(fmt.Sprintf(":%s", port)); err != nil {
		log.Println(err)
	}
}

func Sync() error {
	db, err := factory.DB().(*xorm.Engine)
	if !err {
		return errors.New("cannot sync models with DB")
	}

	CheckErr(db.Sync(new(models.Brand)))
	CheckErr(db.Sync(new(models.Category)))
	CheckErr(db.Sync(new(models.Food)))
	CheckErr(db.Sync(new(models.Nutrient)))
	/*
	CheckErr(db.Sync(new(models.BrandFood)))
	CheckErr(db.Sync(new(models.CategoryFood)))
	CheckErr(db.Sync(new(models.FoodNutrient)))
	*/

	return nil
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getCurPath() (string, error) {
	curPath, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return curPath, nil
}
