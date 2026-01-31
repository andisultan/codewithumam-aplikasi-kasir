package repositories

import (
	"aplikasi-kasir/models"
	"database/sql"
	"errors"
)

type ProductCategoryRepository struct {
	db *sql.DB
}

func NewProductCategoryRepository(db *sql.DB) *ProductCategoryRepository {
	return &ProductCategoryRepository{db: db}
}

func (repo *ProductCategoryRepository) GetAll() ([]models.ProductCategory, error) {
	query := "SELECT id, name FROM product_categories"
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.ProductCategory, 0)
	for rows.Next() {
		var p models.ProductCategory
		err := rows.Scan(&p.ID, &p.Name)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (repo *ProductCategoryRepository) Create(productCategory *models.ProductCategory) error {
	query := "INSERT INTO product_categories (name) VALUES ($1) RETURNING id"
	err := repo.db.QueryRow(query, productCategory.Name).Scan(&productCategory.ID)
	return err
}

// GetByID - ambil product categories by ID
func (repo *ProductCategoryRepository) GetByID(id int) (*models.ProductCategory, error) {
	query := "SELECT id, name FROM product_categories WHERE id = $1"

	var p models.ProductCategory
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name)
	if err == sql.ErrNoRows {
		return nil, errors.New("produk category tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (repo *ProductCategoryRepository) Update(product *models.ProductCategory) error {
	query := "UPDATE product_categories SET name = $1 WHERE id = $2"
	result, err := repo.db.Exec(query, product.Name, product.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("produk category tidak ditemukan")
	}

	return nil
}

func (repo *ProductCategoryRepository) Delete(id int) error {
	query := "DELETE FROM product_categories WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("produk category tidak ditemukan")
	}

	return err
}
