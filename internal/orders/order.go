// Package orders содержит структуры данных и ошибки для работы с заказами.
package orders

import "errors"

type Order struct {
	ID          int    `json:"id" db:"id"`
	Customer    string `json:"customer" db:"customer"`
	Address     string `json:"address" db:"address"`
	IsDelivered bool   `json:"is_delivered" db:"is_delivered"`
}

func (o *Order) MarkDelivered() {
	o.IsDelivered = true
}

var ErrOrderNotFound = errors.New("orders not found")
