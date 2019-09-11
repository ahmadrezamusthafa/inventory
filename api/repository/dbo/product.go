package dbo

import "github.com/jinzhu/gorm"

type Product struct {
	gorm.Model
	SKU  string `gorm:"type:varchar(300);NOT NULL;UNIQUE"`
	Name string `gorm:"type:varchar(200);NOT NULL"`
}

func (Product) TableName() string {
	return "product"
}
