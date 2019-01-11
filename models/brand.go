package models

import "time"

type Brand struct {
	Id			int64		`json:"id"`
	Name		string  	`json:"name"`
	ImgSrc		string		`json:"img_url"`
	CategoryId	int64		`json:"category_id"`
	CreatedAt	time.Time	`json:"created_at" xorm:"created"`
	DeletedAt 	time.Time	`json:"deleted_at" xorm:"deleted"`
}

func (b Brand) Get(id int) error {
	return nil
}
