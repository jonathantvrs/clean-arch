package model

import (
	"time"
)

type Order struct {
	ID          int       `json:"id" db:"id"`
	ProductName string    `json:"product_name" db:"product_name"`
	Quantity    int       `json:"quantity" db:"quantity"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
