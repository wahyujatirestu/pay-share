package model

import (
	"time"

	"github.com/google/uuid"
)

type Transactions struct {
	ID             uuid.UUID 	`db:"id" json:"id"`
	BillDate       *time.Time 	`db:"bill_date" json:"billDate,omitempty"`
	EntryDate      *time.Time 	`db:"entry_date" json:"entryDate,omitempty"`
	FinishDate     *time.Time 	`db:"finish_date" json:"finishDate,omitempty"`
	CustomerId     uuid.UUID 	`db:"customer_id" json:"customerId"`
	EmployeeId     *uuid.UUID 	`db:"employee_id" json:"employeeId,omitempty"`
	TotalAmount    float64   	`db:"total_amount" json:"totalAmount"`
	PaymentStatus  string    	`db:"payment_status" json:"paymentStatus"`
	InvoiceNumber  *string   	`db:"invoice_number" json:"invoiceNumber,omitempty"`
	PaymentMethod  *string   	`db:"payment_method" json:"paymentMethod,omitempty"`
	PaymentURL     *string   	`db:"payment_url" json:"paymentUrl,omitempty"`
	Status         string    	`db:"status" json:"status"`
	Notes          *string   	`db:"notes" json:"notes,omitempty"`
	Created_At      time.Time 	`db:"created_at" json:"createdAt"`
	Updated_At      time.Time 	`db:"updated_at" json:"updatedAt"`
}
