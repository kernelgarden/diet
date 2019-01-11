package models

import "time"

type Nutrient struct {
	Id				int64		`json:"id" xorm:"pk autoincr"`
	FoodId			int64		`xorm:"index"`
	Carbohydrate	float32		`json:"carbohydrate"`
	Protein			float32		`json:"protein"`
	SaturatedFat	float32 	`json:"saturated_fat"`
	UnSaturatedFat	float32 	`json:"unsaturated_fat"`
	TransFat		float32		`json:"trans_fat"`
	PerWeight		int32		`json:"per_weight"`
	Calorie			int64		`json:"calorie"`
	CreatedAt		time.Time	`json:"created_at" xorm:"created"`
	DeletedAt 		time.Time	`json:"deleted_at" xorm:"deleted"`
}