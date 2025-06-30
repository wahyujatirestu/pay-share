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

type TransactionRepository interface {
	Create(transaction *model.Transactions) error
	GetAll(filters map[string]interface{}) ([]*model.Transactions, error)
	GetById(id string) (*model.Transactions, error)
	Update(transaction *model.Transactions) error
	Delete(id string) error
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ts *model.Transactions) error {
	ts.ID = uuid.New()
	now := time.Now()
	ts.Created_At = now
	ts.Updated_At = now

	_, err := r.db.Exec(`INSERT INTO transaction (id, bill_date, entry_date, finish_date, customer_id, employee_id, total_amount, payment_status, invoice_number, payment_method, payment_url, status, notes, created_at, updated_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`)

	return err
}

func (r *transactionRepository) GetAll(filters map[string]interface{}) ([]*model.Transactions, error) {
	query := `SELECT id, bill_date, entry_date, finish_date, customer_id, employee_id, total_amount, payment_status, invoice_number, payment_method, payment_url, status, notes, created_at, updated_at FROM transactions`
	var conditions []string
	var args []interface{}
	i := 1

	for k, v := range filters{
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
		return nil, err
	}

	defer rows.Close()

	var transactions []*model.Transactions
	for rows.Next(){
		ts := model.Transactions{}
		err := rows.Scan(
			&ts.ID,
			&ts.BillDate,
			&ts.EntryDate,
			&ts.FinishDate,
			&ts.CustomerId,
			&ts.EmployeeId,
			&ts.TotalAmount,
			&ts.PaymentStatus,
			&ts.InvoiceNumber,
			&ts.PaymentMethod,
			&ts.PaymentURL,
			&ts.Status,
			&ts.Notes,
			&ts.Created_At,
			&ts.Updated_At,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &ts)
	}

	return transactions, nil
}


func (r *transactionRepository) GetById(id string) (*model.Transactions, error) {
	query := `SELECT id, bill_date, entry_date, finish_date, customer_id, employee_id, total_amount, payment_status, invoice_number, payment_method, payment_url, status, notes, created_at, updated_at FROM transactions WHERE id=$1`
	row := r.db.QueryRow(query, id)

	var tr model.Transactions
	err := row.Scan(
		&tr.ID,
		&tr.BillDate,
		&tr.EntryDate,
		&tr.FinishDate,
		&tr.CustomerId,
		&tr.EmployeeId,
		&tr.TotalAmount,
		&tr.PaymentStatus,
		&tr.InvoiceNumber,
		&tr.PaymentMethod,
		&tr.PaymentURL,
		&tr.Status,
		&tr.Notes,
		&tr.Created_At,
		&tr.Updated_At,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &tr, err
}

func (r *transactionRepository) Update(ts *model.Transactions) error {
	res, err := r.db.Exec(`UPDATE transactions SET bill_date=$1, entry_date=$2, finish_date=$3, customer_id=$4, employee_id=$5, total_amount=$6, payment_status=$7, invoice_number=$8, payment_method=$9, payment_url=$10, status=$11, note=$12 WHERE id=$13`, ts.BillDate, ts.EntryDate, ts.FinishDate, ts.CustomerId, ts.EmployeeId, ts.TotalAmount, ts.PaymentStatus, ts.InvoiceNumber, ts.PaymentMethod, ts.PaymentURL, ts.Status, ts.Notes, ts.ID)
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

func (r *transactionRepository) Delete(id string) error {
	res, err := r.db.Exec(`DELETE FROM transactions WHERE id=$1`, id)
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