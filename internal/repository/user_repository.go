package repository

import (
	"database/sql"
	"fmt"
	"soulstreet/internal/model"
)

type ProductRepository interface {
	Create(product *model.Product) error
	GetByID(id int) (*model.Product, error)
	GetAll() ([]model.Product, error)
	Delete(id int) error
	GetByName(name string) ([]*model.Product, error)
}

type productRepositoryDB struct {
	db *sql.DB
}


func NewProductRepositoryDB(db *sql.DB) ProductRepository {
	return &productRepositoryDB{db: db}
}

func (r *productRepositoryDB) Create(product *model.Product) error {
	query := "INSERT INTO products (name, price, images, sizes) VALUES (?,?,?,?)"
	_, err := r.db.Exec(query, product.Name, product.Price, product.Images, product.Sizes, product.IsAvaliable)
	if err != nil {
		return err
	}
	return nil
}

func (r *productRepositoryDB) GetByID(id int) (*model.Product, error) {
	var product model.Product
	query := "SELECT id, name, price, images, sizes, is_avaliable FROM products WHERE id = ?"
	err := r.db.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.Price, &product.Images, &product.Sizes, &product.IsAvaliable)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("produto não encontrado")
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepositoryDB) GetAll() ([]model.Product, error) {
	var products []model.Product
	rows, err := r.db.Query("SELECT id, name, price, images, sizes, is_avaliable FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var product model.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Images, &product.Sizes, &product.IsAvaliable); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (r *productRepositoryDB) Delete(id int) error {
	query := "DELETE FROM products WHERE id = ?"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *productRepositoryDB) GetByName(name string) ([]*model.Product, error) {
	var products []*model.Product
	
	query := "SELECT id, name, price, images, sizes, is_avaliable FROM products WHERE name LIKE ?"
	
	rows, err := r.db.Query(query, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product model.Product
	
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Images, &product.Sizes, &product.IsAvaliable)
		if err != nil {
			return nil, err
		}
	
		products = append(products, &product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	
	if len(products) == 0 {
		return nil, fmt.Errorf("nenhum produto encontrado com o nome '%s'", name)
	}

	return products, nil
}
