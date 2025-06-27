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

type UserService interface {
	Register(user *model.User, confirmPassword string) error
	Login(email, password string) (*model.User, error)
	GetCustomerById(id string) (*model.User, error)
	GetAllCustomer(filters map[string]interface{})([]*model.User, error)
	Update(customer *model.User) error
	Delete(id string) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
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

func (s *userService) Register(u *model.User, confirmPassword string) error {
	if u.Password != confirmPassword {
		return errors.New("password and confirm password not match")
	}

	existing, err := s.repo.GetByEmail(u.Email)
	if err != nil {
		return err
	}

	if existing != nil {
		return errors.New("email already exist")
	}

	hash, err := generatePasswordHash(u.Password)
	if err != nil {
		return err
	}
	u.Password = hash

	return s.repo.Create(u)
}

func (s *userService) Login(email, password string) (*model.User, error) {
	u, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, errors.New("email not found")
	}

	valid, err := verifyPasswordHash(u.Password, password)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("invalid password")
	}

	return  u, nil
}

func (s *userService) GetCustomerById(id string) (*model.User, error) {
	return s.repo.GetById(id)
}

func (s *userService) GetAllCustomer(filters map[string]interface{}) ([]*model.User, error) {
	return s.repo.GetAll(filters)
}

func (s *userService) Update(u *model.User) error {
	if u.Password != "" {
		hash, err := generatePasswordHash(u.Password)
		if err != nil {
			return err
		}
		u.Password = hash
	}
	return s.repo.Update(u)
}

func (s *userService) Delete(id string) error {
	return s.repo.Delete(id)
}