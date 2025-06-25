package model

import (
	"time"
)

type TransactionDetails struct {
	ID            string    `db:"id" json:"id"`
	TransactionId string    `db:"transaction_id" json:"transactionId"`
	ProductId     string    `db:"product_id" json:"productId"`
	ProductPrice  int       `db:"product_price" json:"productPrice"`
	Qty           int       `db:"qty" json:"qty"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time `db:"updated_at" json:"updatedAt"`
}
