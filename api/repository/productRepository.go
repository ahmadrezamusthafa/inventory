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

func (repository *ProductRepository) GetAvailableStock(productID int) (stock int, err error) {

	var (
		incomming IncommingDetailRepository
		outgoing  OutgoingRepository

		incommingTotal int
		outgoingTotal  int
	)

	incommingTotal, err = incomming.GetIncommingTotalByProduct(productID)
	if err != nil {
		return 0, err
	}

	outgoingTotal, err = outgoing.GetOutgoingTotalByProduct(productID)
	if err != nil {
		return 0, err
	}

	return incommingTotal - outgoingTotal, nil
}

func (repository *ProductRepository) IsProductAvailable(id int) bool {
	var row dbo.Product
	return !repository.databaseORM.First(&row, "id = ?", id).RecordNotFound()
}

func (repository *ProductRepository) IsSKUAvailable(sku string) bool {
	var row dbo.Product
	return !repository.databaseORM.First(&row, "sku = ?", sku).RecordNotFound()
}

func (repository *ProductRepository) IsNameAvailable(name string) bool {
	var row dbo.Product
	return !repository.databaseORM.First(&row, "name = ?", name).RecordNotFound()
}

func (repository *ProductRepository) Create(product dbo.Product) error {
	db := repository.databaseORM.Create(&product)
	return db.Error
}
