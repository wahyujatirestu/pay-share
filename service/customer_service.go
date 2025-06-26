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

type CustomerService interface {
	Register(customer *model.Customer, confirmPassword string) error
	Login(email, password string) (*model.Customer, error)
	GetCustomerById(id string) (*model.Customer, error)
	GetAllCustomer(filters map[string]interface{})([]*model.Customer, error)
	Update(customer *model.Customer) error
	Delete(id string) error
}

type customerService struct {
	repo repository.CustomerRepository
}

func NewCustomerService(repo repository.CustomerRepository) CustomerService {
	return &customerService{repo: repo}
}

func generatePasswordHash1(password string) (string, error) {
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

func verifyPasswordHash1(encodedHash, password string) (bool, error) {
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

func (s *customerService) Register(c *model.Customer, confirmPassword string) error {
	if c.Password != confirmPassword {
		return errors.New("password and confirm password not match")
	}

	existing, err := s.repo.GetByEmail(c.Email)
	if err != nil {
		return err
	}

	if existing != nil {
		return errors.New("email already exist")
	}

	hash, err := generatePasswordHash1(c.Password)
	if err != nil {
		return err
	}
	c.Password = hash

	return s.repo.Create(c)
}

func (s *customerService) Login(email, password string) (*model.Customer, error) {
	c, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, errors.New("email not found")
	}

	valid, err := verifyPasswordHash1(c.Password, password)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("invalid password")
	}

	return  c, nil
}

func (s *customerService) GetCustomerById(id string) (*model.Customer, error) {
	return s.repo.GetById(id)
}

func (s *customerService) GetAllCustomer(filters map[string]interface{}) ([]*model.Customer, error) {
	return s.repo.GetAll(filters)
}

func (s *customerService) Update(c *model.Customer) error {
	if c.Password != "" {
		hash, err := generatePasswordHash1(c.Password)
		if err != nil {
			return err
		}
		c.Password = hash
	}
	return s.repo.Update(c)
}

func (s *customerService) Delete(id string) error {
	return s.repo.Delete(id)
}