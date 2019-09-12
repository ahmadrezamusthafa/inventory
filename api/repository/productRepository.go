package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/rezamusthafa/inventory/api/repository/dbo"
	"github.com/rezamusthafa/inventory/api/repository/types"
	"github.com/rezamusthafa/inventory/api/services/inputs"
	"github.com/rezamusthafa/inventory/util"
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

func (repository *ProductRepository) GetProduct(filter inputs.Filter) ([]types.Product, error) {

	var products []types.Product
	rows, err := repository.databaseORM.Raw(
		`select id, sku, name, created_at from product where created_at >= ? and created_at <= ? order by created_at desc`, filter.StartDate, filter.EndDate).Rows()
	if err != nil {
		return []types.Product{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var product types.Product
		err := rows.Scan(&product.ID, &product.SKU, &product.Name, &product.CreatedAt)
		if err != nil {
			return []types.Product{}, err
		}

		product.StrCreatedAt = product.CreatedAt.Format(util.Timestamp)
		incommingTotal, err := getIncommingTotalByProduct(product.ID, repository.databaseORM)
		if err != nil {
			return []types.Product{}, err
		}

		outgoingTotal, err := getOutgoingTotalByProduct(product.ID, repository.databaseORM)
		if err != nil {
			return []types.Product{}, err
		}

		product.Stock = incommingTotal - outgoingTotal
		if err != nil {
			return []types.Product{}, err
		}

		products = append(products, product)
	}

	return products, nil
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
