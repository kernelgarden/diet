package models

import (
	"github.com/go-xorm/xorm"
	"github.com/kernelgarden/diet/factory"
	"time"
)

type Nutrient struct {
	Id             int64     `json:"Id" xorm:"pk autoincr"`
	FoodId         int64     `json:"FoodId" xorm:"index"`
	Carbohydrate   float32   `json:"Carbohydrate"`
	Protein        float32   `json:"Protein"`
	SaturatedFat   float32   `json:"SaturatedFat"`
	UnSaturatedFat float32   `json:"UnSaturatedFat"`
	TransFat       float32   `json:"TransFat"`
	PerUnit 	   int32     `json:"PerUnit"`
	Calorie        int64     `json:"Calorie"`
	Unit		   int32	 `json:"Unit"`
	CreatedAt      time.Time `json:"-" xorm:"created"`
	DeletedAt      time.Time `json:"-" xorm:"deleted"`
}

func (n *Nutrient) Create() (int64, error) {
	return factory.DB().Insert(n)
}

func (n *Nutrient) CreateWithSes(session *xorm.Session) (int64, error) {
	return session.Insert(n)
}

func (Nutrient) Get(id int64) (*Nutrient, error) {
	var n Nutrient
	if has, err := factory.DB().ID(id).Get(&n); err != nil {
		return &n, err
	} else if !has {
		return nil, nil
	}

	return &n, nil
}

func (Nutrient) GetAll(offset, limit int) ([]*Nutrient, error) {
	// TODO: Increase performance via goroutine
	nutrients := make([]*Nutrient, 0)

	if err := factory.DB().Limit(limit, offset).Find(&nutrients); err != nil {
		return nil, err
	}

	return nutrients, nil
}

func (n *Nutrient) Update() error {
	_, err := factory.DB().ID(n.Id).Update(n)
	return err
}

func (n *Nutrient) UpdateWithSes(session *xorm.Session) error {
	_, err := session.ID(n.Id).Update(n)
	return err
}

func (Nutrient) Delete(id int64) error {
	_, err := factory.DB().ID(id).Delete(&Nutrient{})
	return err
}
