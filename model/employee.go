package model

import (
	"time"

	"github.com/google/uuid"
)

type Employee struct {
	ID 			uuid.UUID		`db:"id" json:"id"`
	Name 		string			`db:"name" json:"name"`
	Email 		string			`db:"email" json:"email"`
	Phone 		string			`db:"phone" json:"phone"`
	Password	string			`db:"password" json:"-"`
	Address 	string			`db:"address" json:"address"`
	Created_At 	time.Time		`db:"created_at" json:"created_at"`
	Updated_At 	time.Time		`db:"updated_at" json:"updated_at"`
}