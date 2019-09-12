package dbo

import (
	"github.com/jinzhu/gorm"
)

type IncommingProduct struct {
	gorm.Model
	Product       *Product
	ProductID     int     `gorm:"type:integer REFERENCES product(id) ON DELETE CASCADE ON UPDATE CASCADE"`
	Receipt       *string `gorm:"type:varchar(100);UNIQUE"`
	OrderQty      int     `gorm:"type:integer"`
	PurchasePrice float64 `gorm:"type:float8;NOT NULL"`
	TotalPrice    float64 `gorm:"type:float8;NOT NULL"`
}

func (IncommingProduct) TableName() string {
	return "incomming_product"
}
