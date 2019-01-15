package factory

import (
	"context"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/kernelgarden/diet/constant"
	"fmt"
	"runtime"
	"github.com/kernelgarden/diet/models"
)

var db *xorm.Engine

func InitDB(driver, connection string) (*xorm.Engine, error) {
	db, err := xorm.NewEngine(driver, connection)
	if err != nil {
		return nil, err
	}

	if driver == "sqlite3" {
		runtime.GOMAXPROCS(1)
	}

	db.Sync(new(models.Brand))
	db.Sync(new(models.Category))
	db.Sync(new(models.Food))
	db.Sync(new(models.Nutrient))

	return db, nil
}

func DB() xorm.Interface {
	if db == nil {
		panic("cannot use DB")
	}

	return db
}

func Transaction(queryExec func() error) error {
	session := db.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return err
	}

	if err := queryExec(); err != nil {
		return err
	}

	return session.Commit()
}