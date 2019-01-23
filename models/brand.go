package models

import (
	"github.com/kernelgarden/diet/factory"
	"time"
)

type Brand struct {
	Id         int64     `json:"Id"`
	Name       string    `json:"Name"`
	ImgSrc     string    `json:"ImgSrc"`
	CategoryId int64     `json:"CategoryId"`
	CreatedAt  time.Time `json:"-" xorm:"created"`
	DeletedAt  time.Time `json:"-" xorm:"deleted"`
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

func (b Brand) Foods() ([]*Food, error) {
	foods := make([]*Food, 0)
	if err := factory.DB().Where("brand_id = ?", b.Id).Find(&foods); err != nil {
		return nil, err
	}

	return foods, nil
}
