package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/wahyujatirestu/payshare/model"
	"github.com/wahyujatirestu/payshare/repository"

	"golang.org/x/crypto/argon2"
)

type EmployeeService interface {
	Register(employee *model.Employee, confirmPassword string) error
	Login(email, password string) (*model.Employee, error)
	GetByID(id string) (*model.Employee, error)
	GetAll(filters map[string]interface{}) ([]*model.Employee, error)
	Update(employee *model.Employee) error
	Delete(id string) error
}

type employeeService struct {
	repo repository.EmployeeRepository
}

func NewEmployeeService(repo repository.EmployeeRepository) EmployeeService {
	return &employeeService{repo: repo}
}

func generatePasswordHash(password string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := strings.Join([]string{"argon2id", "v=19", b64Salt, b64Hash}, "$")

	return encoded, nil

}

func verifyPasswordHash(encodedHash, password string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 4 {
		return false, errors.New("invalid password hash")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[2])
	if err != nil {
		return false, err
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return false, err
	}

	computedHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	return string(computedHash) == string(hash), nil
}

func (s *employeeService) Register(employee *model.Employee, confirmPassword string) error {
	if employee.Password != confirmPassword {
		return errors.New("password and confirm password do not match")
	}

	existing, err := s.repo.GetByEmail(employee.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("email already registered")
	}

	hash, err := generatePasswordHash(employee.Password)
	if err != nil {
		return err
	}
	employee.Password = hash

	return s.repo.Create(employee)
}

func (s *employeeService) Login(email, password string) (*model.Employee, error) {
	emp, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if emp == nil {
		return nil, errors.New("email not found")
	}

	valid, err := verifyPasswordHash(emp.Password, password)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("invalid password")
	}

	return emp, nil
}

func (s *employeeService) GetByID(id string) (*model.Employee, error) {
	return s.repo.GetById(id)
}

func (s *employeeService) GetAll(filters map[string]interface{}) ([]*model.Employee, error) {
	return s.repo.GetAll(filters)
}

func (s *employeeService) Update(employee *model.Employee) error {
	if employee.Password != "" {
		hash, err := generatePasswordHash(employee.Password)
		if err != nil {
			return err
		}
		employee.Password = hash
	}
	return s.repo.Update(employee)
}

func (s *employeeService) Delete(id string) error {
	return s.repo.Delete(id)
}


