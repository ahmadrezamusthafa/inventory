package repository

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/rezamusthafa/inventory/api/repository/dbo"
	"github.com/rezamusthafa/inventory/api/repository/types"
	"github.com/rezamusthafa/inventory/api/services/inputs"
	"github.com/rezamusthafa/inventory/util"
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

func (repository *OutgoingRepository) GetOutgoingProduct(filter inputs.Filter) ([]types.OutgoingProduct, error) {

	var outgoingProducts []types.OutgoingProduct
	rows, err := repository.databaseORM.Raw(
		`select op.id, p.id as product_id, p.sku, p.name, op.order_qty, op.order_code, op.selling_price, op.total_price, p.created_at from outgoing_product op, product p
             where op.product_id = p.id and op.created_at >= ? and op.created_at <= ? order by op.created_at desc`, filter.StartDate, filter.EndDate).Rows()
	if err != nil {
		return []types.OutgoingProduct{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var outgoingProduct types.OutgoingProduct
		err := rows.Scan(&outgoingProduct.ID, &outgoingProduct.ProductID, &outgoingProduct.SKU, &outgoingProduct.Name, &outgoingProduct.OrderQty, &outgoingProduct.OrderCode, &outgoingProduct.SellingPrice, &outgoingProduct.TotalPrice, &outgoingProduct.CreatedAt)
		if err != nil {
			return []types.OutgoingProduct{}, err
		}

		outgoingProduct.StrCreatedAt = outgoingProduct.CreatedAt.Format(util.Timestamp)
		outgoingProduct.Note = fmt.Sprintf("Pesanan %s", outgoingProduct.OrderCode)
		outgoingProducts = append(outgoingProducts, outgoingProduct)
	}

	return outgoingProducts, nil
}

func (repository *OutgoingRepository) GetSalesReport(filter inputs.Filter) ([]types.SalesReport, error) {

	var outgoingProducts []types.SalesReport
	rows, err := repository.databaseORM.Raw(
		`select op.id, p.id as product_id, p.sku, p.name, op.order_qty, op.order_code, op.selling_price, op.total_price, p.created_at from outgoing_product op, product p
             where op.product_id = p.id and op.created_at >= ? and op.created_at <= ? order by op.created_at desc`, filter.StartDate, filter.EndDate).Rows()
	if err != nil {
		return []types.SalesReport{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var outgoingProduct types.SalesReport
		err := rows.Scan(&outgoingProduct.ID, &outgoingProduct.ProductID, &outgoingProduct.SKU, &outgoingProduct.Name, &outgoingProduct.OrderQty, &outgoingProduct.OrderCode, &outgoingProduct.SellingPrice, &outgoingProduct.TotalPrice, &outgoingProduct.CreatedAt)
		if err != nil {
			return []types.SalesReport{}, err
		}

		outgoingProduct.StrCreatedAt = outgoingProduct.CreatedAt.Format(util.Timestamp)

		avgPurchasePrice, err := getAveragePurchasePriceByProduct(outgoingProduct.ProductID, repository.databaseORM)
		if err != nil {
			return []types.SalesReport{}, err
		}

		totalPurchasePrice := avgPurchasePrice * float64(outgoingProduct.OrderQty)
		profit := outgoingProduct.TotalPrice - totalPurchasePrice
		outgoingProduct.AvgPurchasePrice = avgPurchasePrice
		outgoingProduct.Profit = profit

		outgoingProducts = append(outgoingProducts, outgoingProduct)
	}

	return outgoingProducts, nil
}

func (repository *OutgoingRepository) IsOrderCodeAndProductAvailable(orderCode string, productID int) bool {
	var row dbo.OutgoingProduct
	return !repository.databaseORM.First(&row, "order_code = ? and product_id = ?", orderCode, productID).RecordNotFound()
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
	rows, err := repository.databaseORM.Raw("select coalesce(sum(order_qty),0) as total from outgoing_product where product_id = ?", productID).Rows()
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

func getOutgoingTotalByProduct(productID int, db *gorm.DB) (int, error) {

	var total int
	rows, err := db.Raw("select coalesce(sum(order_qty),0) as total from outgoing_product where product_id = ?", productID).Rows()
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
