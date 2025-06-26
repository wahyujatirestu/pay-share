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

type CustomerRepository interface {
	Create(customer *model.Customer) error
	GetAll(filters map[string]interface{}) ([]*model.Customer, error)
	GetById(id string)(*model.Customer, error)
	GetByEmail(email string)(*model.Customer, error)
	Update(customer *model.Customer) error
	Delete(id string) error
}

type customerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(c *model.Customer) error {
	c.ID = uuid.New()
	now := time.Now()
	c.Created_At = now
	c.Updated_At = now

	_, err := r.db.Exec(`INSERT INTO customers (id, name , email, phone, password, address, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, c.ID, c.Name, c.Email, c.Phone, c.Password, c.Created_At, c.Updated_At)
	
	if err != nil {
		return err
	}
	return nil
}

func (r *customerRepository) GetAll(filters map[string]interface{}) ([]*model.Customer, error) {
	query := `SELECT id, name, email, phone, password, address, created_at, updated_at FROM customers`
	var conditions []string
	var args []interface{}
	i := 1

	for k, v := range filters{
		conditions = append(conditions, fmt.Sprintf("%s ILIKE %d", k, i))
		args = append(args, "%"+fmt.Sprintf("%v", v)+"%")
		i++
	}

	if len(conditions) > 0 {
		query += "WHERE " + strings.Join(conditions, " AND ")
	}

	query += "ORDER BY created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return  nil, err
	}
	defer rows.Close()

	var customers []*model.Customer
	for rows.Next(){
		var c model.Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Password, &c.Address, &c.Created_At, &c.Updated_At); err != nil {
			return nil, err
		}
		customers = append(customers, &c)
	}

	return  customers, nil
}

func (r *customerRepository) GetById(id string) (*model.Customer, error) {
	row := r.db.QueryRow(`SELECT id, name, email, phone, password, address, created_at, updated_at FROM customers WHERE id = $1`, id)

	var c model.Customer
	if err := row.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Password, &c.Address, &c.Created_At, &c.Updated_At); err != nil {
		return  nil, err
	}

	return &c, nil
}

func (r *customerRepository) GetByEmail(email string) (*model.Customer, error) {
	row := r.db.QueryRow(`SELECT id, name, email, phone, password, address, created_at, updated_at FROM customers WHERE email = $1`, email)

	var c model.Customer
	if err := row.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Password, &c.Address, &c.Created_At, &c.Updated_At); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *customerRepository) Update(c *model.Customer) error {
	result, err := r.db.Exec(`UPDATE customers SET name=$1, email=$2, phone=$3, password=$4, address=$5, updated_at=$6 WHERE id=$7`, c.Name, c.Email, c.Phone, c.Password, c.Address, c.Updated_At, c.ID)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("no rows Updated")
	}

	return  nil
}

func (r *customerRepository) Delete(id string) error {
	result, err := r.db.Exec(`DELETE FROM customers WHERE id=$1`, id)
	if err != nil {
		return err
	} 

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("no rows deleted")
	}
	return nil
}