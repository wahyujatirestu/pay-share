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
	GetByUsername(username string)(*model.User, error)
	GetEmailUsername(identifier string)(*model.User, error)
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

	_, err := r.db.Exec(`INSERT INTO users (
		id, name, email, username, phone, password, address, role, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		u.ID, u.Name, u.Email, u.Username, u.Phone, u.Password, u.Address, u.Role, u.Created_At, u.Updated_At)

	return err
}


func (r *userRepository) GetAll(filters map[string]interface{}) ([]*model.User, error) {
	query := `SELECT id, name, email, username, phone, password, address, role,  created_at, updated_at FROM users`
	var conditions []string
	var args []interface{}
	i := 1

	for k, v := range filters{
		conditions = append(conditions, fmt.Sprintf("%s ILIKE $%d", k, i))
		args = append(args, "%"+fmt.Sprintf("%v", v)+"%")
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

	var users []*model.User
	for rows.Next(){
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Username, &u.Phone, &u.Password, &u.Address, &u.Role, &u.Created_At, &u.Updated_At); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return  users, nil
}

func (r *userRepository) GetById(id string) (*model.User, error) {
	row := r.db.QueryRow(`SELECT id, name, email, username, phone, password, address, role, created_at, updated_at FROM users WHERE id = $1`, id)

	var u model.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Username, &u.Phone, &u.Password, &u.Address, &u.Role, &u.Created_At, &u.Updated_At); err != nil {
		return  nil, err
	}

	return &u, nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	row := r.db.QueryRow(`SELECT id, name, email, username, phone, password, address, role, created_at, updated_at FROM users WHERE email = $1`, email)

	var u model.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Username, &u.Phone, &u.Password, &u.Address, &u.Role, &u.Created_At, &u.Updated_At); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // data tidak ditemukan, bukan error fatal
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	row := r.db.QueryRow(`SELECT id, name, email, username, phone, password, address, role, created_at, updated_at FROM users WHERE username = $1`, username)

	var u model.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Username, &u.Phone, &u.Password, &u.Address, &u.Role, &u.Created_At, &u.Updated_At); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // data tidak ditemukan, bukan error fatal
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetEmailUsername(identifier string) (*model.User, error) {
	var u model.User
	if err := r.db.QueryRow(`SELECT id, email, username, password, role FROM users WHERE email = $1 OR username = $1`, identifier).
		Scan(&u.ID, &u.Email, &u.Username, &u.Password, &u.Role); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) Update(u *model.User) error {
	result, err := r.db.Exec(`UPDATE users SET name=$1, email=$2, username=$3 phone=$4, password=$5, address=$6, role=$7, updated_at=$8 WHERE id=$9`, u.Name, u.Email, u.Username, u.Phone, u.Password, u.Address, u.Role, u.Updated_At, u.ID)

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