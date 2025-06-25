package model

import (
	"time"
)

type Transactions struct {
	ID            string    `db:"id" json:"id"`
	BillDate      time.Time `db:"bill_date" json:"billDate"`
	EntryDate     time.Time `db:"entry_date" json:"entryDate"`
	FinishDate    time.Time `db:"finish_date" json:"finishDate"`
	CustomerId    string    `db:"customer_id" json:"customerId"`
	EmployeeId    string    `db:"employee_id" json:"employeeId"`
	TotalAmount   float64   `db:"total_amount" json:"totalAmount"`
	PaymentStatus string    `db:"payment_status" json:"paymentStatus"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time `db:"updated_at" json:"updatedAt"`
}
