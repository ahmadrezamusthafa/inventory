package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rezamusthafa/inventory/api/configuration"
	"github.com/rezamusthafa/inventory/api/services"
	"github.com/rezamusthafa/inventory/util"
	"github.com/rs/cors"
	"net/http"
	"path"
	"runtime"
)

type Server struct {
	configuration    *configuration.Configuration
	productService   *services.ProductService
	incommingService *services.IncommingService
	outgoingService  *services.OutgoingService
	reportService    *services.ReportService
	migrationService *services.MigrationService
}

func NewServer(
	config *configuration.Configuration,
	productSvc *services.ProductService,
	incommingSvc *services.IncommingService,
	outgoingSvc *services.OutgoingService,
	reportSvc *services.ReportService,
	migrationSvc *services.MigrationService) *Server {

	return &Server{
		configuration:    config,
		productService:   productSvc,
		incommingService: incommingSvc,
		outgoingService:  outgoingSvc,
		reportService:    reportSvc,
		migrationService: migrationSvc,
	}
}

func (s *Server) NewFrontEndRouter() *mux.Router {

	_, runningFile, _, _ := runtime.Caller(1)
	frontendPath := path.Join(path.Dir(runningFile), "../web")

	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(frontendPath))))

	return router
}

func (s *Server) NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/product/create", s.productService.CreateProduct).Methods("POST")
	router.HandleFunc("/product/get", s.productService.GetProduct).Methods("GET")
	router.HandleFunc("/incomming/create", s.incommingService.CreateIncommingProduct).Methods("POST")
	router.HandleFunc("/incomming/get", s.incommingService.GetIncommingProduct).Methods("GET")
	router.HandleFunc("/outgoing/create", s.outgoingService.CreateOutgoingProduct).Methods("POST")
	router.HandleFunc("/outgoing/get", s.outgoingService.GetOutgoingProduct).Methods("GET")
	router.HandleFunc("/report/product/get", s.reportService.GetReportValueOfProduct).Methods("GET")
	router.HandleFunc("/report/sales/get", s.reportService.GetSalesReport).Methods("GET")
	router.HandleFunc("/report/product/export", s.reportService.ExportReportValueOfProduct).Methods("GET")
	router.HandleFunc("/report/sales/export", s.reportService.ExportSalesReport).Methods("GET")
	router.HandleFunc("/migrate/product", s.migrationService.MigrateProductFromSheet).Methods("GET")

	return router
}

func (s *Server) RunFrontEnd() {
	var port = util.ExtractServerAddressPort(s.configuration.App.FrontEndAddress)
	fmt.Println("Starting WEB at http://localhost:" + port + "/")

	router := s.NewFrontEndRouter()
	corsMiddleware := cors.New(cors.Options{
		AllowedMethods: []string{"OPTIONS", "GET", "POST", "PUT", "DELETE"},
		Debug:          false,
	})

	err := http.ListenAndServe(":"+port, corsMiddleware.Handler(router))
	if err != nil {
		fmt.Errorf("ListenAndServe Front End Error: ", err)
	}
}

func (s *Server) Run() {
	var port = util.ExtractServerAddressPort(s.configuration.App.BackEndAddress)
	fmt.Println("Starting API at http://localhost:" + port + "/")

	router := s.NewRouter()
	corsMiddleware := cors.New(cors.Options{
		AllowedMethods: []string{"OPTIONS", "GET", "POST", "PUT", "DELETE"},
		Debug:          false,
	})

	go s.RunFrontEnd()
	err := http.ListenAndServe(":"+port, corsMiddleware.Handler(router))
	if err != nil {
		fmt.Errorf("ListenAndServe Error: ", err)
	}
}
