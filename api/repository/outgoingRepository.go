package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/rezamusthafa/inventory/api/repository/dbo"
)

type OutgoingRepository struct {
	databaseORM *gorm.DB
}

func NewOutgoingRepository(databaseORM *gorm.DB) *OutgoingRepository {
	return &OutgoingRepository{databaseORM: databaseORM}
}

func (repository *OutgoingRepository) Database() *gorm.DB {
	return repository.databaseORM
}

func (repository *OutgoingRepository) GetAll() []dbo.OutgoingProduct {
	var outgoingProducts []dbo.OutgoingProduct
	repository.databaseORM.Debug().Find(&outgoingProducts)
	return outgoingProducts
}

func (repository *OutgoingRepository) IsOrderCodeAvailable(orderCode string) bool {
	var row dbo.OutgoingProduct
	return !repository.databaseORM.First(&row, "order_code = ?", orderCode).RecordNotFound()
}

func (repository *OutgoingRepository) Create(outgoingProduct dbo.OutgoingProduct) error {
	db := repository.databaseORM.Create(&outgoingProduct)
	return db.Error
}

func (repository *OutgoingRepository) GetOutgoingTotalByProduct(productID int) (int, error) {

	var total int
	rows, err := repository.databaseORM.Raw("select sum(order_qty) as total from outgoing_product where product_id = ?", productID).Rows()
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&total)
		if err != nil {
			return 0, err
		}

		break
	}

	return total, nil
}
