package services

import (
	"github.com/rezamusthafa/inventory/api/configuration"
	"github.com/rezamusthafa/inventory/api/repository"
	"github.com/rezamusthafa/inventory/api/response"
	"net/http"
)

type ReportService struct {
	configuration             *configuration.Configuration
	productRepository         *repository.ProductRepository
	incommingRepository       *repository.IncommingRepository
	incommingDetailRepository *repository.IncommingDetailRepository
	outgoingRepository        *repository.OutgoingRepository
}

func NewReportService(
	config *configuration.Configuration,
	productRepo *repository.ProductRepository,
	incommingRepo *repository.IncommingRepository,
	incommingDetailRepo *repository.IncommingDetailRepository,
	outgoingRepo *repository.OutgoingRepository) *ReportService {

	return &ReportService{
		configuration:             config,
		productRepository:         productRepo,
		incommingRepository:       incommingRepo,
		incommingDetailRepository: incommingDetailRepo,
		outgoingRepository:        outgoingRepo,
	}
}

func (service *ReportService) GetReportValueOfProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues := r.URL.Query()
	filterParam, err := validateRequest(queryValues)
	if err != nil {
		response.WriteError(err.Error(), w)
		return
	}

	products, err := service.productRepository.GetProductReport(filterParam)
	if err != nil {
		response.WriteError("Failed to get product value report", w)
		return
	}

	response.WriteSuccess(products, w)

	return
}

func (service *ReportService) GetSalesReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryValues := r.URL.Query()
	filterParam, err := validateRequest(queryValues)
	if err != nil {
		response.WriteError(err.Error(), w)
		return
	}

	products, err := service.outgoingRepository.GetSalesReport(filterParam)
	if err != nil {
		response.WriteError("Failed to get sales report", w)
		return
	}

	response.WriteSuccess(products, w)

	return
}