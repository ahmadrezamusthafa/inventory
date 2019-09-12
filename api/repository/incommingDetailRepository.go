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

func (repository *IncommingDetailRepository) Create(incommingProductDetail dbo.IncommingProductDetail) error {
	db := repository.databaseORM.Create(&incommingProductDetail)
	return db.Error
}

func (repository *IncommingDetailRepository) CreateWithTx(tx *gorm.DB, incommingProductDetail dbo.IncommingProductDetail) error {
	db := tx.Create(&incommingProductDetail)
	return db.Error
}

func (repository *IncommingDetailRepository) GetIncommingTotalByProduct(productID int) (int, error) {

	var total int
	rows, err := repository.databaseORM.Raw("select coalesce(sum(accepted_qty),0) as total from incomming_product_detail ipd, incomming_product ip where ip.id=ipd.incomming_product_id and ip.product_id = ?", productID).Rows()
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

func getIncommingTotalByProduct(productID int, db *gorm.DB) (int, error) {

	var total int
	rows, err := db.Raw("select coalesce(sum(accepted_qty),0) as total from incomming_product_detail ipd, incomming_product ip where ip.id=ipd.incomming_product_id and ip.product_id = ?", productID).Rows()
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

func getAveragePurchasePriceByProduct(productID int, db *gorm.DB) (float64, error) {

	var price float64
	rows, err := db.Raw("select coalesce(avg(purchase_price),0) as total from incomming_product where product_id = ?", productID).Rows()
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&price)
		if err != nil {
			return 0, err
		}

		break
	}

	return price, nil
}
