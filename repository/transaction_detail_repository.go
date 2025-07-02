package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wahyujatirestu/payshare/model"
)

type TransactionDetailsRespository interface {
	Create(detail *model.TransactionDetails) error
	GetById(id string) (*model.TransactionDetails, error)
	GetByTransactionId(transactionId string)([]*model.TransactionDetails, error)
	GetAll(filters map[string]interface{})([]*model.TransactionDetails, error)
	Update(detail *model.TransactionDetails) error
	Delete(id string) error
}

type transactionDetailsRepository struct {
	db *sql.DB
}

func NewTransactionDetailsRepository(db *sql.DB) TransactionDetailsRespository {
	return &transactionDetailsRepository{db: db}
}

func (r *transactionDetailsRepository) Create(detail *model.TransactionDetails) error {
	detail.ID = uuid.New()
	now := time.Now()
	detail.CreatedAt = now
	detail.UpdatedAt = now
	
	_, err := r.db.Exec(`INSERT INTO transaction_details (id, transaction_id, product_id, product_price, qty, discount_amount, subtotal, status, notes, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, detail.ID, detail.TransactionId, detail.ProductId, detail.ProductPrice, detail.Qty, detail.DiscountAmount, detail.Subtotal, detail.Status, detail.Notes, detail.CreatedAt, detail.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *transactionDetailsRepository) GetById(id string) (*model.TransactionDetails, error) {
	row := r.db.QueryRow(`SELECT id, bill_date, entry_date, finish_date, customer_id, employee_id, total_amount, payment_status, invoice_number, payment_method, payment_url, status, notes, created_at, updated_at FROM transaction_details WHERE id = $1`, id)

	var detail model.TransactionDetails
	err := row.Scan(
		&detail.ID,
		&detail.TransactionId,
		&detail.ProductId,
		&detail.ProductPrice,
		&detail.Qty,
		&detail.DiscountAmount,
		&detail.Subtotal,
		&detail.Status,
		&detail.Notes,
		&detail.CreatedAt,
		&detail.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &detail, err
}

func (r *transactionDetailsRepository) GetByTransactionId(transactionId string) ([]*model.TransactionDetails, error) {
	rows, err := r.db.Query(`SELECT id, transaction_id, product_id, product_price, qty, discount_amount, subtotal, status, notes, created_at, updated_at FROM transaction_details WHERE transaction_id = $1`, transactionId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var details []*model.TransactionDetails
	for rows.Next() {
		var d model.TransactionDetails
		err := rows.Scan(
			&d.ID,
			&d.TransactionId,
			&d.ProductId,
			&d.ProductPrice,
			&d.Qty,
			&d.DiscountAmount,
			&d.Subtotal,
			&d.Status,
			&d.Notes,
			&d.CreatedAt,
			&d.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		details = append(details, &d)
	}

	return details, nil
}

func (r *transactionDetailsRepository) GetAll(filters map[string]interface{}) ([]*model.TransactionDetails, error) {
	query := `SELECT id, transaction_id, product_id, product_price, qty, discount_amount, subtotal, status, notes, created_at, updated_at FROM transaction_details`
	var conditions []string
	var args []interface{}
	i := 1

	for k, v := range filters {
		conditions = append(conditions, fmt.Sprintf("%s=$%d", k, i))
		args = append(args, v)
		i++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return  nil, err
	}

	defer rows.Close()

	var details []*model.TransactionDetails
	for rows.Next(){
		var d model.TransactionDetails
		err := rows.Scan(
			&d.ID,
			&d.TransactionId,
			&d.ProductId,
			&d.ProductPrice,
			&d.Qty,
			&d.DiscountAmount,
			&d.Subtotal,
			&d.Status,
			&d.Notes,
			&d.CreatedAt,
			&d.UpdatedAt,
		)

		if err != nil {
			return  nil, err
		}

		details = append(details, &d)
	}
	return details, nil
}

func (r *transactionDetailsRepository) Update(detail *model.TransactionDetails) error {
	res, err := r.db.Exec(`UPDATE transaction_details SET transaction_id=$1, product_id=$2, product_price=$3, qty=$4, discount_amount=$5, subtotal=$6, status=$7, notes=$8 WHERE id=$9`, detail.TransactionId, detail.ProductId, detail.ProductPrice, detail.Qty, detail.DiscountAmount, detail.Subtotal, detail.Status, detail.Notes, detail.ID)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func (r *transactionDetailsRepository) Delete(id string) error {
	res, err := r.db.Exec(`DELETE FROM transaction_details WHERE id = $1`, id)
	if err != nil {
		return err
	}

	row, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if row == 0 {
		return errors.New("no rows affected")
	}

	return nil
}