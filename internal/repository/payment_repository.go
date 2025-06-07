package repository

import (
	"database/sql"
	
	"fmt"
	"soulstreet/internal/model"
)

type PaymentRepository interface {
	Save(payment *model.Payment) error
	GetProductInfo(productID int) (*model.Product, error)
}

type PaymentRepositoryDB struct {
	db *sql.DB
}

func NewPaymentRepositoryDB(db *sql.DB) PaymentRepository {
	return &PaymentRepositoryDB{db: db}
}

func (r *PaymentRepositoryDB) Save(payment *model.Payment) error {
	fmt.Printf("Saving payment %s in database\n", payment.ID)
	
	res, err := r.db.Exec(`
			INSERT INTO address (neighborhood, street, number, state, country)
			VALUES (?, ?, ?, ?, ?)`,
		payment.Address.Neighborhood,
		payment.Address.Street,
		payment.Address.Number,
		payment.Address.State,
		payment.Address.Country,
	)
	if err != nil {
		return fmt.Errorf("failed to insert address: %w", err)
	}

	addressID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get address id: %w", err)
	}

	_, err = r.db.Exec(`
			INSERT INTO payment (product_id, email, telephone, cpf, address_id)
			VALUES (?, ?, ?, ?, ?)`,
		payment.ProductID,
		payment.Email,
		payment.Telephone,
		payment.CPF,
		addressID,
	)
	if err != nil {
		return fmt.Errorf("failed to insert payment: %w", err)
	}

	fmt.Println("Payment saved successfully.")
	return nil

}

func (r *PaymentRepositoryDB) GetProductInfo(productID int) (*model.Product, error) {
	var product model.Product
	query := "SELECT id, name, price, images, sizes, is_avaliable FROM products WHERE id = ?"
	err := r.db.QueryRow(query, productID).Scan(&product.ID, &product.Name, &product.Price, &product.Images, &product.Sizes, &product.IsAvaliable)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("produto não encontrado")
		}
		return nil, err
	}
	return &product, nil
}
