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
	CreatedAt  time.Time `json:"-" xorm:"created"`
	DeletedAt  time.Time `json:"-" xorm:"deleted"`
}

type FoodJSON struct {
	Food     Food     `json:"food"`
	Nutrient Nutrient `json:"nutrient"`
	Brand    Brand    `json:"brand"`
	Category Category `json:"category"`
}

func (f Food) ToJSON() (*FoodJSON, error) {
	var nutrient Nutrient
	if has, err := factory.DB().Where("food_id = ?", f.Id).Get(&nutrient); err != nil {
		return nil, err
	} else if !has {
		return nil, nil
	}

	// 브랜드는 없는 경우도 있을 수 있다.
	var brand Brand
	if _, err := factory.DB().ID(f.BrandId).Get(&brand); err != nil {
		return nil, err
	}

	var category Category
	if has, err := factory.DB().ID(f.CategoryId).Get(&category); err != nil {
		return nil, err
	} else if !has {
		return nil, nil
	}

	return &FoodJSON{Food: f, Nutrient: nutrient, Brand: brand, Category: category}, nil
}

func (FoodJSON) NewFoodJSON(food Food, nutrient Nutrient, brand Brand, category Category) FoodJSON {
	return FoodJSON{Food: food, Nutrient: nutrient, Brand: brand, Category: category}
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
