package models

import "time"

type Food struct {
	Id         int64     `json:"id" xorm:"pk autoincr"`
	CategoryId int64     `xorm:"index"`
	BrandId    int64     `xorm:"index"`
	Name       string    `json:"name" xorm:"varchar(64)"`
	Weight     float64   `json:"weight"`
	CreatedAt  time.Time `json:"created_at" xorm:"created"`
	DeletedAt  time.Time `json:"deleted_at" xorm:"deleted"`
}

type FoodNutrient struct {
	Food     `xorm:"extends"`
	Nutrient `xorm:"extends"`
}

func (f *Food) Create() error {
	return nil
}

func (f *Food) GetAll() error {
	return nil
}

func (f *Food) Get(id int) error {
	return nil
}

func (f Food) Update() error {
	return nil
}

func (f Food) Delete() error {
	return nil
}
