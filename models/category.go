package models

import (
	"github.com/kernelgarden/diet/factory"
	"time"
)

type Category struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at" xorm:"created"`
	DeletedAt time.Time `json:"deleted_at" xorm:"deleted"`
}

type CategoryFood struct {
	Category `xorm:"extends"`
	Food     `xorm:"extends"`
}

func (CategoryFood) TableName() string {
	return "food"
}

func (c *Category) Create() (int64, error) {
	return factory.DB().Insert(c)
}

func (Category) Get(id int64) (*Category, error) {
	var c Category
	if has, err := factory.DB().ID(id).Get(&c); err != nil {
		return &c, err
	} else if !has {
		return nil, nil
	}

	return &c, nil
}

func (Category) GetAll(offset, limit int) ([]*Category, error) {
	// TODO: Increase performance via goroutine
	categories := make([]*Category, 0)

	if err := factory.DB().Limit(limit, offset).Find(&categories); err != nil {
		return nil, err
	}

	return categories, nil
}

func (c *Category) Update() error {
	_, err := factory.DB().ID(c.Id).Update(c)
	return err
}

func (Category) Delete(id int64) error {
	_, err := factory.DB().ID(id).Delete(&Category{})
	return err
}
