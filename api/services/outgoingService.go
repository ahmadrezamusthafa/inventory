package services

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/rezamusthafa/inventory/api/configuration"
	"github.com/rezamusthafa/inventory/api/repository"
	"github.com/rezamusthafa/inventory/api/repository/dbo"
	"github.com/rezamusthafa/inventory/api/response"
	"github.com/rezamusthafa/inventory/api/response/results"
	"github.com/rezamusthafa/inventory/api/services/inputs"
	"io/ioutil"
	"net/http"
)

type OutgoingService struct {
	configuration             *configuration.Configuration
	productRepository         *repository.ProductRepository
	outgoingRepository        *repository.OutgoingRepository
	incommingDetailRepository *repository.IncommingDetailRepository
}

func NewOutgoingService(
	config *configuration.Configuration,
	productRepo *repository.ProductRepository,
	outgoingRepo *repository.OutgoingRepository,
	incommingDetailRepo *repository.IncommingDetailRepository) *OutgoingService {

	return &OutgoingService{
		configuration:             config,
		productRepository:         productRepo,
		outgoingRepository:        outgoingRepo,
		incommingDetailRepository: incommingDetailRepo,
	}
}

func (service *OutgoingService) GetOutgoingProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues := r.URL.Query()
	filterParam, err := validateRequest(queryValues)
	if err != nil {
		response.WriteError(err.Error(), w)
		return
	}

	outgoingProducts, err := service.outgoingRepository.GetOutgoingProduct(filterParam)
	if err != nil {
		response.WriteError("Failed to get outgoing product", w)
		return
	}

	response.WriteSuccess(outgoingProducts, w)

	return
}

func (service *OutgoingService) CreateOutgoingProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var (
		err            error
		data           []byte
		product        inputs.Outgoing
		availableStock int
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

	if product.ProductID == 0 || product.OrderQty == 0 || product.SellingPrice == 0 || product.OrderCode == "" {
		response.WriteError("Missing request parameter", w)
		return
	}

	availableStock, err = service.GetAvailableStock(product.ProductID)
	if err != nil {
		response.WriteError("Failed to get available stock", w)
		return
	}

	if availableStock < product.OrderQty {
		response.WriteError(fmt.Sprintf("Available product stock is %d", availableStock), w)
		return
	}

	if !service.productRepository.IsProductAvailable(product.ProductID) {
		response.WriteError("Product id is unavailable", w)
		return
	}

	if service.outgoingRepository.IsOrderCodeAndProductAvailable(product.OrderCode, product.ProductID) {
		response.WriteError("Order code and product id is already exist", w)
		return
	}

	err = service.outgoingRepository.Create(dbo.OutgoingProduct{
		ProductID:    product.ProductID,
		OrderCode:    product.OrderCode,
		OrderQty:     product.OrderQty,
		SellingPrice: product.SellingPrice,
		TotalPrice:   (float64(product.OrderQty) * product.SellingPrice),
	})
	if err != nil {
		response.WriteError("Failed to create outgoing product", w)
	}

	var successObj = results.TransactionStatus{Message: "Successfully created outgoing product"}
	response.WriteSuccess(successObj, w)
}

func (service *OutgoingService) GetAvailableStock(productID int) (stock int, err error) {

	var (
		incommingTotal int
		outgoingTotal  int
	)

	incommingTotal, err = service.incommingDetailRepository.GetIncommingTotalByProduct(productID)
	if err != nil {
		return 0, err
	}

	outgoingTotal, err = service.outgoingRepository.GetOutgoingTotalByProduct(productID)
	if err != nil {
		return 0, err
	}

	return (incommingTotal - outgoingTotal), nil
}
