package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/rezamusthafa/inventory/api/repository/dbo"
)

type IncommingDetailRepository struct {
	databaseORM *gorm.DB
}

func NewIncommingDetailRepository(databaseORM *gorm.DB) *IncommingDetailRepository {
	return &IncommingDetailRepository{databaseORM: databaseORM}
}

func (repository *IncommingDetailRepository) Database() *gorm.DB {
	return repository.databaseORM
}

func (repository *IncommingDetailRepository) GetAll() []dbo.IncommingProductDetail {
	var incommingProductDetails []dbo.IncommingProductDetail
	repository.databaseORM.Debug().Find(&incommingProductDetails)
	return incommingProductDetails
}