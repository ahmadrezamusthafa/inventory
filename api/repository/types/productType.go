package types

import "time"

type Product struct {
	ID           int       `json:"id"`
	SKU          string    `json:"sku"`
	Name         string    `json:"name"`
	Stock        int       `json:"stock"`
	CreatedAt    time.Time `json:"created_at"`
	StrCreatedAt string    `json:"str_created_at"`
}
