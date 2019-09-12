package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rezamusthafa/inventory/api/configuration"
	"github.com/rezamusthafa/inventory/api/services"
	"github.com/rezamusthafa/inventory/util"
	"github.com/rs/cors"
	"net/http"
)

type Server struct {
	configuration    *configuration.Configuration
	productService   *services.ProductService
	incommingService *services.IncommingService
	outgoingService  *services.OutgoingService
}

func NewServer(
	config *configuration.Configuration,
	productSvc *services.ProductService,
	incommingSvc *services.IncommingService,
	outgoingSvc *services.OutgoingService) *Server {

	return &Server{
		configuration:    config,
		productService:   productSvc,
		incommingService: incommingSvc,
		outgoingService:  outgoingSvc,
	}
}

func (s *Server) NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/product/create", s.productService.CreateProduct).Methods("POST")
	router.HandleFunc("/product/get", s.productService.GetProduct).Methods("GET")
	router.HandleFunc("/incomming/create", s.incommingService.CreateIncommingProduct).Methods("POST")
	router.HandleFunc("/incomming/get", s.incommingService.GetIncommingProduct).Methods("GET")
	router.HandleFunc("/outgoing/create", s.outgoingService.CreateOutgoingProduct).Methods("POST")
	router.HandleFunc("/outgoing/get", s.outgoingService.GetOutgoingProduct).Methods("GET")

	return router
}

func (s *Server) Run() {
	var port = util.ExtractServerAddressPort(s.configuration.App.BackEndAddress)
	fmt.Println("Starting API at http://localhost:" + port + "/")

	router := s.NewRouter()
	corsMiddleware := cors.New(cors.Options{
		AllowedMethods: []string{"OPTIONS", "GET", "POST", "PUT", "DELETE"},
		Debug:          false,
	})

	err := http.ListenAndServe(":"+port, corsMiddleware.Handler(router))
	if err != nil {
		fmt.Errorf("ListenAndServe Error: ", err)
	}
}
