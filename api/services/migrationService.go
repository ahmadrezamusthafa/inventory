package services

import (
	"fmt"
	"github.com/rezamusthafa/inventory/api/configuration"
	"github.com/rezamusthafa/inventory/api/repository"
	"github.com/rezamusthafa/inventory/api/response"
	"github.com/rezamusthafa/inventory/api/response/results"
	"github.com/rezamusthafa/inventory/api/services/core"
	"net/http"
)

type MigrationService struct {
	configuration     *configuration.Configuration
	productRepository *repository.ProductRepository
}

func NewMigrationService(
	config *configuration.Configuration,
	productRepo *repository.ProductRepository) *MigrationService {

	return &MigrationService{
		configuration:     config,
		productRepository: productRepo,
	}
}

func (service *MigrationService) MigrateProductFromSheet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	products, err := core.ReadSheetProduct(service.configuration)
	if err != nil {
		fmt.Println(err)
		response.WriteError("Failed to migrate product from sheet", w)
		return
	}

	for _, product := range products {
		if service.productRepository.IsSKUAvailable(product.SKU) {
			continue
		}

		if service.productRepository.IsNameAvailable(product.Name) {
			continue
		}

		err = service.productRepository.Create(product)
		if err != nil {
			fmt.Println("Failed to migrate product", product.Name)
			continue
		}
	}

	var successObj = results.TransactionStatus{Message: "Migrate successfully executed"}
	response.WriteSuccess(successObj, w)
}
