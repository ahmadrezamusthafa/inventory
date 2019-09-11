package api

import (
	"github.com/rezamusthafa/inventory/api/configuration"
	"github.com/rezamusthafa/inventory/api/services"
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
