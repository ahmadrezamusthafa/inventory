package inputs

type Outgoing struct {
	ID           uint    `json:"id"`
	SKU          string  `json:"sku",omitempty`
	ProductID    int     `json:"product_id"`
	OrderCode    string  `json:"order_code"`
	OrderQty     int     `json:"order_qty"`
	SellingPrice float64 `json:"selling_price"`
}
