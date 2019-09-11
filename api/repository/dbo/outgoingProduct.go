package dbo

import "github.com/jinzhu/gorm"

type OutgoingProduct struct {
	gorm.Model
	Product      *Product
	ProductId    int     `gorm:"type:integer REFERENCES product(id) ON DELETE CASCADE ON UPDATE CASCADE"`
	OrderID      string  `gorm:"type:varchar(100);NOT NULL;UNIQUE"`
	OrderQty     int     `gorm:"type:integer"`
	SellingPrice float64 `gorm:"type:float8;NOT NULL"`
	TotalPrice   float64 `gorm:"type:float8;NOT NULL"`
}

func (OutgoingProduct) TableName() string {
	return "outgoing_product"
}
