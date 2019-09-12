package repository

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/rezamusthafa/inventory/api/repository/dbo"
	"github.com/rezamusthafa/inventory/api/repository/types"
	"github.com/rezamusthafa/inventory/api/services/inputs"
	"github.com/rezamusthafa/inventory/util"
	"strings"
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

func (repository *IncommingRepository) GetIncommingProductDetail(incommingProductID int) ([]types.IncommingProductDetail, types.IncommingProductDetail, error) {

	var (
		incommingProductDetails []types.IncommingProductDetail
		summaryDetail           types.IncommingProductDetail
		strb                    strings.Builder
		index                   int
	)
	rows, err := repository.databaseORM.Raw(
		`select id, incomming_product_id, accepted_qty, created_at from incomming_product_detail where incomming_product_id = ?;`, incommingProductID).Rows()
	if err != nil {
		return []types.IncommingProductDetail{}, types.IncommingProductDetail{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var incommingProductDetail types.IncommingProductDetail
		err := rows.Scan(&incommingProductDetail.ID, &incommingProductDetail.IncommingProductID, &incommingProductDetail.AcceptedQty, &incommingProductDetail.CreatedAt)
		if err != nil {
			return []types.IncommingProductDetail{}, types.IncommingProductDetail{}, err
		}

		incommingProductDetail.StrCreatedAt = incommingProductDetail.CreatedAt.Format(util.Timestamp)
		summaryDetail.AcceptedQty += incommingProductDetail.AcceptedQty
		strb.WriteString(fmt.Sprintf("%s%s terima %d", func() string {
			if index > 0 {
				return ", "
			}
			return ""
		}(), incommingProductDetail.CreatedAt.Format(util.Timestamp), incommingProductDetail.AcceptedQty))
		incommingProductDetails = append(incommingProductDetails, incommingProductDetail)

		index++
	}

	summaryDetail.Note = strb.String()
	return incommingProductDetails, summaryDetail, nil
}

func (repository *IncommingRepository) GetIncommingProduct(filter inputs.Filter) ([]types.IncommingProduct, error) {

	var incommingProducts []types.IncommingProduct
	rows, err := repository.databaseORM.Raw(
		`select ip.id, p.id, p.sku, p.name, ip.order_qty, ip.purchase_price, ip.total_price, COALESCE(ip.receipt,'(hilang)') receipt, p.created_at from incomming_product ip, product p
             where ip.product_id = p.id and ip.created_at >= ? and ip.created_at <= ? order by ip.created_at desc`, filter.StartDate, filter.EndDate).Rows()
	if err != nil {
		return []types.IncommingProduct{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var incommingProduct types.IncommingProduct
		err := rows.Scan(&incommingProduct.ID, &incommingProduct.ProductID, &incommingProduct.SKU, &incommingProduct.Name, &incommingProduct.OrderQty, &incommingProduct.PurchasePrice, &incommingProduct.TotalPrice, &incommingProduct.Receipt, &incommingProduct.CreatedAt)
		if err != nil {
			return []types.IncommingProduct{}, err
		}

		detail, summary, err := repository.GetIncommingProductDetail(incommingProduct.ID)
		if err != nil {
			return []types.IncommingProduct{}, err
		}

		incommingProduct.StrCreatedAt = incommingProduct.CreatedAt.Format(util.Timestamp)
		incommingProduct.IncommingDetail = detail
		incommingProduct.Note = summary.Note
		incommingProduct.AcceptedQty = summary.AcceptedQty

		incommingProducts = append(incommingProducts, incommingProduct)
	}

	return incommingProducts, nil
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

func (repository *IncommingRepository) CreateAndReturnIDWithTx(tx *gorm.DB, incommingProduct dbo.IncommingProduct) (uint, error) {
	var row dbo.IncommingProduct
	db := tx.Create(&incommingProduct).Scan(&row)
	if db.Error != nil {
		return 0, db.Error
	}

	return row.ID, nil
}
