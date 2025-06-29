package service

import (
	"errors"

	"github.com/wahyujatirestu/payshare/model"
	"github.com/wahyujatirestu/payshare/repository"
	"github.com/wahyujatirestu/payshare/utils/security"
)

type UserService interface {
	Register(user *model.User, confirmPassword string) error
	GetUserById(id string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	GetEmailUsername(identifier string) (*model.User, error)
	GetAllUser(filters map[string]interface{})([]*model.User, error)
	Update(customer *model.User) error
	Delete(id string) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}


func (s *userService) Register(u *model.User, confirmPassword string) error {
	if u.Password != confirmPassword {
		return errors.New("password and confirm password not match")
	}

	existingEmail, err := s.repo.GetByEmail(u.Email)
	if err != nil {
		return err
	}

	if existingEmail != nil {
		return errors.New("email already exist")
	}

	existingUsername, err := s.repo.GetByUsername(u.Username)
	if err != nil {
		return err
	}

	if existingUsername != nil {
		return errors.New("username already exist")
	}

	hash, err := security.GeneratePasswordHash(u.Password)
	if err != nil {
		return err
	}
	u.Password = hash

	return s.repo.Create(u)
}

func (s *userService) GetUserById(id string) (*model.User, error) {
	return s.repo.GetById(id)
}

func (s *userService) GetUserByEmail(email string) (*model.User, error) {
	return s.repo.GetByEmail(email)
}

func (s *userService) GetEmailUsername(identifier string) (*model.User, error) {
	return s.repo.GetEmailUsername(identifier)
}

func (s *userService) GetAllUser(filters map[string]interface{}) ([]*model.User, error) {
	return s.repo.GetAll(filters)
}

func (s *userService) Update(u *model.User) error {
	if u.Password != "" {
		hash, err := security.GeneratePasswordHash(u.Password)
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