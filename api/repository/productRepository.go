package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/rezamusthafa/inventory/api/repository/dbo"
)

type ProductRepository struct {
	databaseORM *gorm.DB
}

func NewProductRepository(databaseORM *gorm.DB) *ProductRepository {
	return &ProductRepository{databaseORM: databaseORM}
}

func (repository *ProductRepository) Database() *gorm.DB {
	return repository.databaseORM
}

func (repository *ProductRepository) GetAll() []dbo.Product {
	var products []dbo.Product
	repository.databaseORM.Debug().Find(&products)
	return products
}
