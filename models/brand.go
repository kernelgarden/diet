package models

import (
	"github.com/kernelgarden/diet/factory"
	"time"
)

type Brand struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	ImgSrc     string    `json:"img_url"`
	CategoryId int64     `json:"category_id"`
	CreatedAt  time.Time `json:"created_at" xorm:"created"`
	DeletedAt  time.Time `json:"deleted_at" xorm:"deleted"`
}

type BrandFood struct {
	Brand `xorm:"extends"`
	Food  `xorm:"extends"`
}

func (BrandFood) TableName() string {
	return "food"
}

func (b *Brand) Create() (int64, error) {
	return factory.DB().Insert(b)
}

func (Brand) Get(id int64) (*Brand, error) {
	var b Brand
	if has, err := factory.DB().ID(id).Get(&b); err != nil {
		return &b, err
	} else if !has {
		return nil, nil
	}

	return &b, nil
}

func (Brand) GetAll(offset, limit int) ([]*Brand, error) {
	// TODO: Increase performance via goroutine
	brands := make([]*Brand, 0)

	if err := factory.DB().Limit(limit, offset).Find(&brands); err != nil {
		return nil, err
	}

	return brands, nil
}

func (b *Brand) Update() error {
	_, err := factory.DB().ID(b.Id).Update(b)
	return err
}

func (Brand) Delete(id int64) error {
	_, err := factory.DB().ID(id).Delete(&Brand{})
	return err
}

func (Brand) Default() (*Brand, error) {
	var b Brand
	if _, err := factory.DB().ID(0).Get(&b); err != nil {
		return nil, err
	}

	return &b, nil
}
