package services

import (
	"errors"
	"github.com/json-iterator/go"
	"github.com/rezamusthafa/inventory/api/configuration"
	"github.com/rezamusthafa/inventory/api/repository"
	"github.com/rezamusthafa/inventory/api/repository/dbo"
	"github.com/rezamusthafa/inventory/api/response"
	"github.com/rezamusthafa/inventory/api/response/results"
	"github.com/rezamusthafa/inventory/api/services/inputs"
	"io/ioutil"
	"net/http"
	"strings"
)

type IncommingService struct {
	configuration             *configuration.Configuration
	productRepository         *repository.ProductRepository
	incommingRepository       *repository.IncommingRepository
	incommingDetailRepository *repository.IncommingDetailRepository
}

func NewIncommingService(
	config *configuration.Configuration,
	productRepo *repository.ProductRepository,
	incommingRepo *repository.IncommingRepository,
	incommingDetailRepo *repository.IncommingDetailRepository) *IncommingService {

	return &IncommingService{
		configuration:             config,
		productRepository:         productRepo,
		incommingRepository:       incommingRepo,
		incommingDetailRepository: incommingDetailRepo,
	}
}

func (service *IncommingService) GetIncommingProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues := r.URL.Query()
	filterParam, err := validateRequest(queryValues)
	if err != nil {
		response.WriteError(err.Error(), w)
		return
	}

	incommingProducts, err := service.incommingRepository.GetIncommingProduct(filterParam)
	if err != nil {
		response.WriteError("Failed to get incomming product", w)
		return
	}

	response.WriteSuccess(incommingProducts, w)

	return
}

func (service *IncommingService) CreateIncommingProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var (
		err     error
		data    []byte
		product inputs.Incomming
		receipt *string
	)

	data, err = ioutil.ReadAll(r.Body)
	if err != nil {
		response.WriteError("Failed to read body request", w)
		return
	}

	err = jsoniter.Unmarshal([]byte(data), &product)
	if err != nil {
		response.WriteError("Failed unmarshal", w)
		return
	}

	if (product.ProductID == 0 && product.SKU == "") || product.OrderQty == 0 || product.AcceptedQty == 0 || product.PurchasePrice == 0 {
		response.WriteError("Missing request parameter", w)
		return
	}

	if product.ProductID <= 0 {
		product.SKU = strings.ToUpper(product.SKU)
		product.ProductID, err = service.productRepository.GetProductIDBySKU(product.SKU)
	}

	if product.AcceptedQty > product.OrderQty {
		response.WriteError("Accepted quantity is invalid", w)
		return
	}

	if !service.productRepository.IsProductAvailable(product.ProductID) {
		response.WriteError("Product id is unavailable", w)
		return
	}

	if product.Receipt != "" {
		if service.incommingRepository.IsReceiptAvailable(product.Receipt) {
			response.WriteError("Receipt is already exist", w)
			return
		}
		receipt = &product.Receipt
	}

	err = executeInsertIncommingProduct(service, product, receipt)
	if err != nil {
		response.WriteError(err.Error(), w)
		return
	}

	var successObj = results.TransactionStatus{Message: "Successfully created incomming product"}
	response.WriteSuccess(successObj, w)
}

func executeInsertIncommingProduct(service *IncommingService, product inputs.Incomming, receipt *string) error {
	tx := service.incommingRepository.Database().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return errors.New("Transaction initial error")
	}

	incommingID, err := service.incommingRepository.CreateAndReturnIDWithTx(tx, dbo.IncommingProduct{
		ProductID:     product.ProductID,
		Receipt:       receipt,
		OrderQty:      product.OrderQty,
		PurchasePrice: product.PurchasePrice,
		TotalPrice:    (float64(product.OrderQty) * product.PurchasePrice),
	})
	if err != nil {
		tx.Rollback()
		return errors.New("Failed to create incomming product")
	}

	if incommingID == 0 {
		return errors.New("Incomming id is invalid")
	}

	err = service.incommingDetailRepository.CreateWithTx(tx, dbo.IncommingProductDetail{
		IncommingProductID: int(incommingID),
		AcceptedQty:        product.AcceptedQty,
	})
	if err != nil {
		tx.Rollback()
		return errors.New("Failed to create incomming product detail")
	}

	if tx.Commit().Error != nil {
		return errors.New("Transaction commit error")

	}

	return nil
}
