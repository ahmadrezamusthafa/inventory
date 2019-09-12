package types

import "time"

type OutgoingProduct struct {
	ID           int       `json:"id"`
	ProductID    int       `json:"product_id"`
	SKU          string    `json:"sku"`
	Name         string    `json:"name"`
	OrderQty     int       `json:"order_qty"`
	OrderCode    string    `json:"order_code"`
	SellingPrice float64   `json:"selling_price"`
	TotalPrice   float64   `json:"total_price"`
	CreatedAt    time.Time `json:"created_at"`
	StrCreatedAt string    `json:"str_created_at"`
	Note         string    `json:"note"`
}
