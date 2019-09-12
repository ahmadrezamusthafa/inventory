package types

import "time"

type IncommingProduct struct {
	ID              int                      `json:"id"`
	ProductID       int                      `json:"product_id"`
	SKU             string                   `json:"sku"`
	Name            string                   `json:"name"`
	OrderQty        int                      `json:"order_qty"`
	AcceptedQty     int                      `json:"accepted_qty"`
	PurchasePrice   float64                  `json:"purchase_price"`
	TotalPrice      float64                  `json:"total_price"`
	Receipt         string                   `json:"receipt"`
	IncommingDetail []IncommingProductDetail `json:"incomming_detail"`
	CreatedAt       time.Time                `json:"created_at"`
	StrCreatedAt    string                   `json:"str_created_at"`
	Note            string                   `json:"note"`
}

type IncommingProductDetail struct {
	ID                 int       `json:"id"`
	IncommingProductID int       `json:"incomming_product_id"`
	AcceptedQty        int       `json:"accepted_qty"`
	CreatedAt          time.Time `json:"created_at"`
	StrCreatedAt       string    `json:"str_created_at"`
	Note               string    `json:"note,omitempty"`
}
