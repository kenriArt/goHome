// Package orders содержит структуры данных и ошибки для работы с заказами.
package orders

import "errors"

type Order struct {
	ID          int
	Customer    string
	Address     string
	IsDelivered bool
}

func (o *Order) MarkDelivered() {
	o.IsDelivered = true
}

var ErrOrderNotFound = errors.New("orders not found")
