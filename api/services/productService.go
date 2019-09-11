package services

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/rezamusthafa/inventory/api/configuration"
	"github.com/rezamusthafa/inventory/api/repository"
	"github.com/rezamusthafa/inventory/api/repository/dbo"
	"github.com/rezamusthafa/inventory/api/response"
	"github.com/rezamusthafa/inventory/api/response/results"
	"github.com/rezamusthafa/inventory/api/services/core"
	"github.com/rezamusthafa/inventory/api/services/inputs"
	"io/ioutil"
	"net/http"
	"strings"
)

type ProductService struct {
	configuration             *configuration.Configuration
	productRepository         *repository.ProductRepository
	incommingRepository       *repository.IncommingRepository
	incommingDetailRepository *repository.IncommingDetailRepository
	outgoingRepository        *repository.OutgoingRepository
}

func NewProductService(
	config *configuration.Configuration,
	productRepo *repository.ProductRepository,
	incommingRepo *repository.IncommingRepository,
	incommingDetailRepo *repository.IncommingDetailRepository,
	outgoingRepo *repository.OutgoingRepository) *ProductService {

	return &ProductService{
		configuration:             config,
		productRepository:         productRepo,
		incommingRepository:       incommingRepo,
		incommingDetailRepository: incommingDetailRepo,
		outgoingRepository:        outgoingRepo,
	}
}

func (service *ProductService) CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var (
		err     error
		data    []byte
		product inputs.Product
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

	if product.Name == "" || product.Size == "" || product.Color == "" {
		response.WriteError("Missing request parameter", w)
		return
	}

	sku := core.GenerateSKU(product)
	if sku == "" {
		response.WriteError("Failed generate SKU", w)
		return
	}

	if service.productRepository.IsSKUAvailable(sku) {
		response.WriteError("Code is already exist", w)
		return
	}

	name := fmt.Sprintf("%s (%s, %s)", product.Name, strings.ToUpper(product.Size), product.Color)
	if service.productRepository.IsNameAvailable(name) {
		response.WriteError("Name is already exist", w)
		return
	}

	err = service.productRepository.Create(dbo.Product{
		SKU:  sku,
		Name: name,
	})
	if err != nil {
		response.WriteError("Failed create product", w)
	}

	var successObj = results.TransactionStatus{Message: "Successfully created product"}
	response.WriteSuccess(successObj, w)
}
