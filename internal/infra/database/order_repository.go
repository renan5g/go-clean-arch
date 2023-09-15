package database

import (
	"database/sql"

	"github.com/renan5g/go-clean-arch/internal/domain/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(conn *sql.DB) *OrderRepository {
	return &OrderRepository{Db: conn}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	query := `INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)`

	stmt, err := r.Db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) GetTotal() (int, error) {
	var total int
	err := r.Db.QueryRow("SELECT count(*) FROM orders").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *OrderRepository) List() ([]*entity.Order, error) {
	rows, err := r.Db.Query("SELECT id, price, tax, final_price FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []*entity.Order{}

	for rows.Next() {
		var id string
		var price, tax, finalPrice float64
		if err := rows.Scan(&id, &price, &tax, &finalPrice); err != nil {
			return nil, err
		}
		orders = append(orders, &entity.Order{ID: id, Price: price, Tax: tax, FinalPrice: finalPrice})
	}

	return orders, nil
}
