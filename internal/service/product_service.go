package service

import (
	"fmt"
	"soulstreet/internal/model"
	"soulstreet/internal/repository"
)


type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) * ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(product *model.Product) error {
	if err := s.repo.Create(product); err != nil {
		return fmt.Errorf("Erro ao criar produto: %v", err)
	}
	return nil
}

func (s *ProductService) GetProductByID(id int) (*model.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) GetAll() ([]model.Product, error) {
	return s.repo.GetAll()
}


func (s *ProductService) DeleteProduct(id int) error {
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("Erro ao deletar produto: %v", err)
	}
	return nil
}
func (s *ProductService) GetProductByName(name string) ([]*model.Product, error) {
	product, err := s.repo.GetByName(name)
	if err != nil {
		return nil, fmt.Errorf("Erro ao buscar produto por nome: %v", err)
	}
	return product, nil
}