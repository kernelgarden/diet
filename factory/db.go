package factory

import (
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/kernelgarden/diet/constant"
	"runtime"
	"sync"
)

var db *xorm.Engine
var once sync.Once

func InitDB() (*xorm.Engine, error) {
	config := constant.GlobalCtx.Value(constant.CtxConfig)
	if config == nil {
		return nil, errors.New("config file is empty")
	}

	c, isValid := config.(*constant.Config)
	if !isValid {
		return nil, errors.New("config file is broken")
	}

	// Init DB
	var connection string
	if c.Debug {
		if c.Database.Host != "" {
			connection = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s_dev?charset=utf8", c.Database.Username, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Name)
		} else {
			connection = fmt.Sprintf("%s:%s@/%s_dev?charset=utf8", c.Database.Username, c.Database.Password, c.Database.Name)
		}
	} else {
		if c.Database.Host != "" {
			connection = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.Database.Username, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Name)
		} else {
			connection = fmt.Sprintf("%s:%s@/%s?charset=utf8", c.Database.Username, c.Database.Password, c.Database.Name)
		}
	}

	var driver string
	if c.Database.Driver == "" {
		driver = "mysql"
	} else {
		driver = c.Database.Driver
	}

	var err error
	db, err = xorm.NewEngine(driver, connection)
	if err != nil {
		return nil, err
	}

	if driver == "sqlite3" {
		runtime.GOMAXPROCS(1)
	}

	return db, nil
}

func DB() xorm.Interface {
	once.Do(func() {
		_, err := InitDB()
		if err != nil {
			panic("cannot connect to DB")
		}
	})

	return db
}

func Transaction(queryExec func(session *xorm.Session) error) error {
	session := db.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return err
	}

	if err := queryExec(session); err != nil {
		session.Rollback()
		return err
	}

	return session.Commit()
}
