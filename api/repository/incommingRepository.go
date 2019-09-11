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