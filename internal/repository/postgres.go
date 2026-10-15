package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"GOTRAKER/internal/orders"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresOrderRepo struct {
	db *sqlx.DB
}

func NewPostgresOrderRepo(db *sqlx.DB) *PostgresOrderRepo {
	return &PostgresOrderRepo{db: db}
}
func (r *PostgresOrderRepo) Add(o *orders.Order) error {
	query := `INSERT INTO orders (id, customer, address, is_delivered) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, o.ID, o.Customer, o.Address, o.IsDelivered)
	if err != nil {
		log.Printf("❌ Ошибка при INSERT в Postgres для заказа %d: %v", o.ID, err)
		return err
	}
	return nil
}

func (r *PostgresOrderRepo) GetByID(ctx context.Context, id int) (*orders.Order, error) {
	fmt.Println("из базы")
	var o orders.Order
	query := `SELECT id, customer, address, is_delivered FROM orders WHERE id = $1`
	err := r.db.GetContext(ctx, &o, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, orders.ErrOrderNotFound
		}
		return nil, fmt.Errorf("ошибка при чтении заказа %d из Postgres: %w", id, err)
	}
	return &o, err
}

func (r *PostgresOrderRepo) Update(o orders.Order) error {
	query := `UPDATE orders SET customer = $1, address = $2, is_delivered = $3 WHERE id = $4`
	result, err := r.db.Exec(query, o.Customer, o.Address, o.IsDelivered, o.ID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return orders.ErrOrderNotFound
	}
	return nil
}

func (r *PostgresOrderRepo) GetAll() []orders.Order {
	query := `SELECT id, customer, address, is_delivered FROM orders`
	rows, err := r.db.Query(query)
	if err != nil {
		fmt.Println("Ошибка при получении заказов:", err)
		return nil
	}
	defer rows.Close()

	var orders1 []orders.Order
	for rows.Next() {
		var o orders.Order
		err := rows.Scan(&o.ID, &o.Customer, &o.Address, &o.IsDelivered)
		if err == nil {
			orders1 = append(orders1, o)
		}
	}
	return orders1
}
