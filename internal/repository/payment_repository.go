package repository

import (
	"database/sql"
	"fmt"
	"soulstreet/internal/model"
)

type PaymentRepository interface {
	Save(payment *model.Payment) error
}

type PaymentRepositoryDB struct {
	db *sql.DB
}

func NewPaymentRepositoryDB(db *sql.DB) PaymentRepository {
	return &PaymentRepositoryDB{db: db}
}

func (r *PaymentRepositoryDB) Save(payment *model.Payment) error {
	fmt.Println("Saving payment to the database")
	return nil
}

func (r *PaymentRepositoryDB) GetPriceByProductID(productID int) (float64, error) {
	var price float64
	query := "SELECT price FROM products WHERE id = ?"
	err := r.db.QueryRow(query, productID).Scan(&price)
	if err != nil {
		return 0, err
	}
	return price, nil
}