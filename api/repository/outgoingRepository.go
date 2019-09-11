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
