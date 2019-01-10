package models

import "time"

type Food struct {
	Id			int64		`json:"id"`
	Name		string		`json:"name"`
	Weight		float64		`json:"weight"`
	CreatedAt	time.Time	`json:"created_at" xorm:"created"`
	DeletedAt 	time.Time	`json:"deleted_at" xorm:"deleted"`
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

}
