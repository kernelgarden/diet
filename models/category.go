package models

import "time"

type Category struct {
	Id			int64	`json:"id"`
	Name		string	`json:"name"`
	CreatedAt	time.Time	`json:"created_at" xorm:"created"`
	DeletedAt 	time.Time	`json:"deleted_at" xorm:"deleted"`
}

type CategoryFood struct {
	Category	`xorm:"extends"`
	Food		`xorm:"extends"`
}
