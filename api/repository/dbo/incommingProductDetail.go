package dbo

import "github.com/jinzhu/gorm"

type IncommingProductDetail struct {
	gorm.Model
	IncommingProduct   *IncommingProduct
	IncommingProductId int `gorm:"type:integer REFERENCES incomming_product(id) ON DELETE CASCADE ON UPDATE CASCADE"`
	AcceptedQty        int `gorm:"type:integer"`
}

func (IncommingProductDetail) TableName() string {
	return "incomming_product_detail"
}
