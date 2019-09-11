package services

import (
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
	configuration      *configuration.Configuration
	productRepository  *repository.ProductRepository
	outgoingRepository *repository.OutgoingRepository
}

func NewOutgoingService(
	config *configuration.Configuration,
	productRepo *repository.ProductRepository,
	outgoingRepo *repository.OutgoingRepository) *OutgoingService {

	return &OutgoingService{
		configuration:      config,
		productRepository:  productRepo,
		outgoingRepository: outgoingRepo,
	}
}

func (service *OutgoingService) CreateOutgoingProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var (
		err     error
		data    []byte
		product inputs.Outgoing
	)

	data, err = ioutil.ReadAll(r.Body)
	if err != nil {
		response.WriteError("Failed read body request", w)
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

	if !service.productRepository.IsProductAvailable(product.ProductID) {
		response.WriteError("Product id is unavailable", w)
		return
	}

	if service.outgoingRepository.IsOrderCodeAvailable(product.OrderCode) {
		response.WriteError("Order code is already exist", w)
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
		response.WriteError("Failed create outgoing product", w)
	}

	var successObj = results.TransactionStatus{Message: "Successfully created outgoing product"}
	response.WriteSuccess(successObj, w)
}
