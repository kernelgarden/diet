package models

import (
	"github.com/go-xorm/xorm"
	"github.com/kernelgarden/diet/factory"
	"time"
)

type Food struct {
	Id         int64     `json:"id" xorm:"pk autoincr"`
	CategoryId int64     `json:"category_id" xorm:"index"`
	BrandId    int64     `json:"brand_id" xorm:"index"`
	Name       string    `json:"name" xorm:"varchar(64)"`
	Weight     float64   `json:"weight"`
	CreatedAt  time.Time `json:"created_at" xorm:"created"`
	DeletedAt  time.Time `json:"deleted_at" xorm:"deleted"`
}

type FoodNutrient struct {
	Food     `xorm:"extends"`
	Nutrient `xorm:"extends"`
}

func (FoodNutrient) TableName() string {
	return "nutrient"
}

func (f *Food) Create() (int64, error) {
	return factory.DB().Insert(f)
}

func (f *Food) CreateWithSes(session *xorm.Session) (int64, error) {
	return session.Insert(f)
}

func (Food) Get(id int64) (*Food, error) {
	var f Food
	if has, err := factory.DB().ID(id).Get(&f); err != nil {
		return &f, err
	} else if !has {
		return nil, nil
	}

	return &f, nil
}

func (Food) GetAll(offset, limit int) ([]*Food, error) {
	// TODO: Increase performance via goroutine
	foods := make([]*Food, 0)

	if err := factory.DB().Limit(limit, offset).Find(&foods); err != nil {
		return nil, err
	}

	return foods, nil
}

func (f *Food) Update() error {
	_, err := factory.DB().ID(f.Id).Update(f)
	return err
}

func (f *Food) UpdateWithSes(session *xorm.Session) error {
	_, err := session.ID(f.Id).Update(f)
	return err
}

func (Food) Delete(id int64) error {
	_, err := factory.DB().ID(id).Delete(&Food{})
	return err
}
