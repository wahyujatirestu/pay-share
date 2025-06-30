package model

import (
	"time"

	"github.com/google/uuid"
)

type TransactionDetails struct {
	ID             uuid.UUID `db:"id" json:"id"`
	TransactionId  uuid.UUID `db:"transaction_id" json:"transactionId"`
	ProductId      uuid.UUID `db:"product_id" json:"productId"`
	ProductPrice   float64   `db:"product_price" json:"productPrice"`
	Qty            int       `db:"qty" json:"qty"`
	DiscountAmount float64   `db:"discount_amount" json:"discountAmount"`
	Subtotal       float64   `db:"subtotal" json:"subtotal"`
	Status         string    `db:"status" json:"status"`
	Notes          *string   `db:"notes" json:"notes,omitempty"`
	CreatedAt      time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time `db:"updated_at" json:"updatedAt"`
}
