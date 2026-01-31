package services

import (
	"aplikasi-kasir/models"
	"aplikasi-kasir/repositories"
)

type ProductCategoryService struct {
	repo *repositories.ProductCategoryRepository
}

func NewProductCategoryService(repo *repositories.ProductCategoryRepository) *ProductCategoryService {
	return &ProductCategoryService{repo: repo}
}

func (s *ProductCategoryService) GetAll() ([]models.ProductCategory, error) {
	return s.repo.GetAll()
}

func (s *ProductCategoryService) Create(data *models.ProductCategory) error {
	return s.repo.Create(data)
}

func (s *ProductCategoryService) GetByID(id int) (*models.ProductCategory, error) {
	return s.repo.GetByID(id)
}

func (s *ProductCategoryService) Update(productCategory *models.ProductCategory) error {
	return s.repo.Update(productCategory)
}

func (s *ProductCategoryService) Delete(id int) error {
	return s.repo.Delete(id)
}
