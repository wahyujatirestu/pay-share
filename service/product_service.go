package service

import (
	"errors"

	"github.com/wahyujatirestu/payshare/model"
	"github.com/wahyujatirestu/payshare/repository"
)

type ProductService interface {
	Create(product *model.Product) error
	GetById(id string) (*model.Product, error)
	GetAll(filters map[string]interface{}) ([]*model.Product, error)
	Update(product *model.Product) error
	Delete(id string) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) Create(p *model.Product) error {
	if p.Name == "" {
		return  errors.New("Product name is required")
	}
	if p.Price <= 0 {
		return errors.New("Product price must be greater than 0")
	}

	return s.repo.Create(p)
}

func (s *productService) GetById(id string) (*model.Product, error) {
	return s.repo.GetById(id)
}

func (s *productService) GetAll(filters map[string]interface{}) ([]*model.Product, error) {
	return s.repo.GetAll(filters)
}

func (s *productService) Update(p *model.Product) error {
	if p.Name == "" {
		return  errors.New("Product name is required")
	}

	if p.Price <= 0 {
		return errors.New("Product price must be greater than 0")
	}

	return s.repo.Update(p)
}

func (s *productService) Delete(id string) error {
	return s.repo.Delete(id)
}