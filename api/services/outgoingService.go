package services

import (
	"github.com/rezamusthafa/inventory/api/configuration"
	"github.com/rezamusthafa/inventory/api/repository"
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
