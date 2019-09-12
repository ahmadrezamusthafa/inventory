package types

import "time"

type ProductReport struct {
	ID               int       `json:"id"`
	SKU              string    `json:"sku" header:"SKU"`
	Name             string    `json:"name" header:"Nama Item"`
	Stock            int       `json:"stock" header:"Jumlah"`
	AvgPurchasePrice float64   `json:"avg_purchase_price" header:"Rata-Rata Harga Beli"`
	TotalPrice       float64   `json:"total_price" header:"Total"`
	CreatedAt        time.Time `json:"created_at"`
	StrCreatedAt     string    `json:"str_created_at"`
}

type SalesReport struct {
	ID               int       `json:"id"`
	ProductID        int       `json:"product_id"`
	OrderCode        string    `json:"order_code" header:"ID Pesanan"`
	StrCreatedAt     string    `json:"str_created_at" header:"Waktu"`
	SKU              string    `json:"sku" header:"SKU"`
	Name             string    `json:"name" header:"Nama Barang"`
	OrderQty         int       `json:"order_qty" header:"Jumlah"`
	SellingPrice     float64   `json:"selling_price" header:"Harga Jual"`
	TotalPrice       float64   `json:"total_price" header:"Total"`
	AvgPurchasePrice float64   `json:"avg_purchase_price" header:"Harga Beli"`
	Profit           float64   `json:"profit" header:"Laba"`
	CreatedAt        time.Time `json:"created_at"`
}
