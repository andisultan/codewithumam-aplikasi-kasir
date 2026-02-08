package repositories

import (
	"aplikasi-kasir/models"
	"database/sql"
	"errors"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// CREATE TABLE IF NOT EXISTS transactions (
// 	id SERIAL PRIMARY KEY,
// 	total_amount INT NOT NULL,
// 	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// );

// CREATE TABLE IF NOT EXISTS transaction_details (
// 	id SERIAL PRIMARY KEY,
// 	transaction_id INT REFERENCES transactions(id) ON DELETE CASCADE,
// 	product_id INT REFERENCES products(id),
// 	product_name VARCHAR(200) NOT NULL,
// 	product_price INT NOT NULL,
// 	quantity INT NOT NULL,
// 	subtotal INT NOT NULL,
// );

func (repo *ProductRepository) GetAll(name string) ([]models.Product, error) {
	query := `
		SELECT
			p.id,
			p.name,
			p.price,
			p.stock,
			p.category_id,
			c.id,
			c.name
		FROM products p
		LEFT JOIN product_categories c
			ON c.id = p.category_id
	`

	var args []any
	if name != "" {
		query += " WHERE p.name ILIKE $1"
		args = append(args, "%"+name+"%")
	}

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []models.Product{}

	for rows.Next() {
		var p models.Product

		var categoryID sql.NullInt64
		var catID sql.NullInt64
		var catName sql.NullString

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Price,
			&p.Stock,
			&categoryID,
			&catID,
			&catName,
		)
		if err != nil {
			return nil, err
		}

		// set category_id
		if categoryID.Valid {
			id := int(categoryID.Int64)
			p.CategoryID = &id
		}

		// set category object (jika ada)
		if catID.Valid {
			id := int(catID.Int64)
			p.Category = &models.ProductCategory{
				ID:   id,
				Name: catName.String,
			}
		}

		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (repo *ProductRepository) Create(product *models.Product) error {
	query := "INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	err := repo.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
	return err
}

func (repo *ProductRepository) GetByID(id int) (*models.Product, error) {
	query := `
		SELECT
			p.id,
			p.name,
			p.price,
			p.stock,
			p.category_id,
			c.id,
			c.name
		FROM products p
		LEFT JOIN product_categories c
			ON c.id = p.category_id
		WHERE p.id = $1
	`

	var p models.Product

	var categoryID sql.NullInt64
	var catID sql.NullInt64
	var catName sql.NullString

	err := repo.db.QueryRow(query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Price,
		&p.Stock,
		&categoryID,
		&catID,
		&catName,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("produk tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	// set category_id
	if categoryID.Valid {
		cid := int(categoryID.Int64)
		p.CategoryID = &cid
	}

	// set category object
	if catID.Valid {
		cid := int(catID.Int64)
		p.Category = &models.ProductCategory{
			ID:   cid,
			Name: catName.String,
		}
	}

	return &p, nil
}

func (repo *ProductRepository) Update(product *models.Product) error {
	query := "UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5"
	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID, product.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}

	return nil
}

func (repo *ProductRepository) Delete(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}

	return err
}
