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

type ProductRepository interface {
	Create(product *model.Product) error
	GetById(id string) (*model.Product, error)
	GetAll(filters map[string]interface{})([]*model.Product, error)
	Update(product *model.Product) error
	Delete(id string) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(p *model.Product) error {
	p.ID = uuid.New()
	now := time.Now()
	p.Created_At = now
	p.Updated_At = now

	_, err := r.db.Exec(`INSERT INTO products (id, name, description, price, unit, created_at, updated_at) VALUES ($1, $2, $3,$4, $5, $6, $7)`, p.ID, p.Name, p.Description, p.Price, p.Unit, p.Created_At, p.Updated_At)
	return  err
}

func (r *productRepository) GetById(id string) (*model.Product, error) {
	var p model.Product
	row := r.db.QueryRow(`SELECT id, name, description, price, unit, created_at, updated_at FROM products WHERE id = $1`, id).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Unit, &p.Created_At, &p.Updated_At)
	if row == sql.ErrNoRows {
		return nil, row
	}
	return &p, nil
}

func (r *productRepository) GetAll(filters map[string]interface{}) ([]*model.Product, error) {
	baseQuery := `SELECT id, name, description, price, unit, created_at, updated_at FROM products`
	var conditions []string
	var args []interface{}
	i := 1

	for k, v := range filters {
		switch k {
		case "price_min":
			conditions = append(conditions, fmt.Sprintf("price >= $%d", i))
			args = append(args, v)
		case "price_max":
			conditions = append(conditions, fmt.Sprintf("price <= $%d", i))
			args = append(args, v)
		default: 
			conditions = append(conditions, fmt.Sprintf("%s ILIKE $%d", k, i))
			args = append(args, "%"+fmt.Sprintf("%v", v)+"%")
		}
		i++
	}

	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}
	baseQuery += " ORDER BY created_at DESC"

	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []*model.Product
	for rows.Next(){
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Unit, &p.Created_At, &p.Updated_At); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}

	return products, nil
}

func (r *productRepository) Update(p *model.Product) error {
	res, err := r.db.Exec(`UPDATE products SET name=$1, description=$2, price=$3, unit=$4 WHERE id=$5`, p.Name, p.Description, p.Price, p.Unit, p.ID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("No rows updated")
	}
	return nil
}

func (r *productRepository) Delete(id string) error {
	res, err := r.db.Exec(`DELETE FROM product WHERE id=$1`, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("No rows deleted")
	}
	return nil
}