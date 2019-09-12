package types

import "time"

type ProductReport struct {
	ID               int       `json:"id"`
	SKU              string    `json:"sku"`
	Name             string    `json:"name"`
	Stock            int       `json:"stock"`
	AvgPurchasePrice float64   `json:"avg_purchase_price"`
	TotalPrice       float64   `json:"total_price"`
	CreatedAt        time.Time `json:"created_at"`
	StrCreatedAt     string    `json:"str_created_at"`
}

type SalesReport struct {
	ID               int       `json:"id"`
	ProductID        int       `json:"product_id"`
	SKU              string    `json:"sku"`
	Name             string    `json:"name"`
	OrderQty         int       `json:"order_qty"`
	OrderCode        string    `json:"order_code"`
	SellingPrice     float64   `json:"selling_price"`
	TotalPrice       float64   `json:"total_price"`
	AvgPurchasePrice float64   `json:"avg_purchase_price"`
	Profit           float64   `json:"profit"`
	CreatedAt        time.Time `json:"created_at"`
	StrCreatedAt     string    `json:"str_created_at"`
}
