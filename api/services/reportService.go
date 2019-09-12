package services

import (
	"fmt"
	"github.com/rezamusthafa/inventory/api/configuration"
	"github.com/rezamusthafa/inventory/api/repository"
	"github.com/rezamusthafa/inventory/api/response"
	"github.com/rezamusthafa/inventory/api/response/results"
	"github.com/rezamusthafa/inventory/api/services/core"
	"github.com/rezamusthafa/inventory/util"
	"net/http"
	"time"
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

func (service *ReportService) ExportReportValueOfProduct(w http.ResponseWriter, r *http.Request) {
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

	var (
		totalProduct int
		totalPrice   float64
	)

	for _, product := range products {
		totalProduct += product.Stock
		totalPrice += product.TotalPrice
	}

	title := [][]string{}
	title = append(title, []string{"LAPORAN NILAI BARANG"})
	title = append(title, []string{""})
	title = append(title, []string{"Tanggal Cetak", time.Now().Format(util.DateOnly)})
	title = append(title, []string{"Jumlah SKU", fmt.Sprintf("%d", len(products))})
	title = append(title, []string{"Jumlah Total Barang", fmt.Sprintf("%d", totalProduct)})
	title = append(title, []string{"Total Nilai", fmt.Sprintf("%.f", totalPrice)})

	err = core.ExportDataToCSV(title, fmt.Sprintf("[Laporan] Nilai Barang %s.csv", time.Now().Format(util.DateOnlySimple)), products)
	if err != nil {
		response.WriteError("Failed to export product value report", w)
		return
	}

	var successObj = results.TransactionStatus{Message: "Successfully export product value report"}
	response.WriteSuccess(successObj, w)

	return
}

func (service *ReportService) ExportSalesReport(w http.ResponseWriter, r *http.Request) {
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

	var (
		totalOmzet   float64
		totalGross   float64
		totalProduct int
	)

	for _, product := range products {
		totalOmzet += product.TotalPrice
		totalGross += product.Profit
		totalProduct += product.OrderQty
	}

	title := [][]string{}
	title = append(title, []string{"LAPORAN PENJUALAN"})
	title = append(title, []string{""})
	title = append(title, []string{"Tanggal Cetak", time.Now().Format(util.DateOnly)})
	title = append(title, []string{"Tanggal", fmt.Sprintf("%s sampai %s", filterParam.StartDate, filterParam.EndDate)})
	title = append(title, []string{"Total Omzet", fmt.Sprintf("%.f", totalOmzet)})
	title = append(title, []string{"Total Laba Kotor", fmt.Sprintf("%.f", totalGross)})
	title = append(title, []string{"Total Penjualan", fmt.Sprintf("%d", len(products))})
	title = append(title, []string{"Total Barang", fmt.Sprintf("%d", totalProduct)})

	err = core.ExportDataToCSV(title, fmt.Sprintf("[Laporan] Penjualan %s.csv", time.Now().Format(util.DateOnlySimple)), products)
	if err != nil {
		response.WriteError("Failed to export sales report", w)
		return
	}

	var successObj = results.TransactionStatus{Message: "Successfully export sales report"}
	response.WriteSuccess(successObj, w)

	return
}
