package models

import (
	"github.com/go-xorm/xorm"
	"github.com/kernelgarden/diet/factory"
	"time"
)

type Nutrient struct {
	Id             int64     `json:"id" xorm:"pk autoincr"`
	FoodId         int64     `json:"food_id" xorm:"index"`
	Carbohydrate   float32   `json:"carbohydrate"`
	Protein        float32   `json:"protein"`
	SaturatedFat   float32   `json:"saturated_fat"`
	UnSaturatedFat float32   `json:"unsaturated_fat"`
	TransFat       float32   `json:"trans_fat"`
	PerWeight      int32     `json:"per_weight"`
	Calorie        int64     `json:"calorie"`
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
