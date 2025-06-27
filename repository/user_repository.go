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

type UserRepository interface {
	Create(user *model.User) error
	GetAll(filters map[string]interface{}) ([]*model.User, error)
	GetById(id string)(*model.User, error)
	GetByEmail(email string)(*model.User, error)
	Update(customer *model.User) error
	Delete(id string) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(u *model.User) error {
	u.ID = uuid.New()
	now := time.Now()
	u.Created_At = now
	u.Updated_At = now

	_, err := r.db.Exec(`INSERT INTO users (id, name , email, phone, password, address, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`, u.ID, u.Name, u.Email, u.Phone, u.Password, u.Role, u.Created_At, u.Updated_At)
	
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetAll(filters map[string]interface{}) ([]*model.User, error) {
	query := `SELECT id, name, email, phone, password, address, role,  created_at, updated_at FROM users`
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

	var users []*model.User
	for rows.Next(){
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.Password, &u.Address, &u.Role, &u.Created_At, &u.Updated_At); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return  users, nil
}

func (r *userRepository) GetById(id string) (*model.User, error) {
	row := r.db.QueryRow(`SELECT id, name, email, phone, password, address, role, created_at, updated_at FROM users WHERE id = $1`, id)

	var u model.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.Password, &u.Address, &u.Role, &u.Created_At, &u.Updated_At); err != nil {
		return  nil, err
	}

	return &u, nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	row := r.db.QueryRow(`SELECT id, name, email, phone, password, address, role, created_at, updated_at FROM users WHERE email = $1`, email)

	var u model.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.Password, &u.Address, &u.Role, &u.Created_At, &u.Updated_At); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) Update(u *model.User) error {
	result, err := r.db.Exec(`UPDATE users SET name=$1, email=$2, phone=$3, password=$4, role=$5 address=$6, updated_at=$7 WHERE id=$8`, u.Name, u.Email, u.Phone, u.Password, u.Address, u.Updated_At, u.ID)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("no rows Updated")
	}

	return  nil
}

func (r *userRepository) Delete(id string) error {
	result, err := r.db.Exec(`DELETE FROM users WHERE id=$1`, id)
	if err != nil {
		return err
	} 

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("no rows deleted")
	}
	return nil
}