package model

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID 				uuid.UUID		`db:"id" json:"id"`
	Name 			string			`db:"name" json:"name"`
	Description		string			`db:"description" json:"description"`
	Price 			float64			`db:"price" json:"price"`
	Unit 			string			`db:"unit" json:"unit"`
	Created_At		time.Time		`db:"created_at" json:"created_at"`
	Updated_At   	time.Time 		`db:"updated_at" json:"updated_at"`
}