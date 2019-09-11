package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/rezamusthafa/inventory/api/repository/dbo"
)

type IncommingRepository struct {
	databaseORM *gorm.DB
}

func NewIncommingRepository(databaseORM *gorm.DB) *IncommingRepository {
	return &IncommingRepository{databaseORM: databaseORM}
}

func (repository *IncommingRepository) Database() *gorm.DB {
	return repository.databaseORM
}

func (repository *IncommingRepository) GetAll() []dbo.IncommingProduct {
	var incommingProducts []dbo.IncommingProduct
	repository.databaseORM.Debug().Find(&incommingProducts)
	return incommingProducts
}

func (repository *IncommingRepository) IsReceiptAvailable(receipt string) bool {
	var row dbo.IncommingProduct
	return !repository.databaseORM.First(&row, "receipt = ?", receipt).RecordNotFound()
}

func (repository *IncommingRepository) Create(incommingProduct dbo.IncommingProduct) error {
	db := repository.databaseORM.Create(&incommingProduct)
	return db.Error
}

func (repository *IncommingRepository) CreateAndReturnID(incommingProduct dbo.IncommingProduct) (uint, error) {
	var row dbo.IncommingProduct
	db := repository.databaseORM.Create(&incommingProduct).Scan(&row)
	if db.Error != nil {
		return 0, db.Error
	}

	return row.ID, nil
}
