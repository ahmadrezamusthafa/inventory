package inputs

type Incomming struct {
	ID            uint    `json:"id"`
	ProductID     int     `json:"product_id"`
	Receipt       string  `json:"receipt"`
	OrderQty      int     `json:"order_qty"`
	AcceptedQty   int     `json:"accepted_qty"`
	PurchasePrice float64 `json:"purchase_price"`
}

type Filter struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
}
