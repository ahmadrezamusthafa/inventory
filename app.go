package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/rezamusthafa/inventory/api/configuration"
	"github.com/rezamusthafa/inventory/api/repository"
	"github.com/rezamusthafa/inventory/api/repository/dbo"
	"github.com/rezamusthafa/inventory/api/services"
)

var (
	config                    *configuration.Configuration
	productRepository         *repository.ProductRepository
	incommingRepository       *repository.IncommingRepository
	incommingDetailRepository *repository.IncommingDetailRepository
	outgoingRepository        *repository.OutgoingRepository
)

var (
	productSvc   *services.ProductService
	incommingSvc *services.IncommingService
	outgoingSvc  *services.OutgoingService
)

func connectToDatabase(config *configuration.Configuration) (*gorm.DB, error) {
	return gorm.Open("sqlite3", config.ConnectionString.Path)
}

func initRepository(db *gorm.DB) {
	productRepository = repository.NewProductRepository(db)
	incommingRepository = repository.NewIncommingRepository(db)
	incommingDetailRepository = repository.NewIncommingDetailRepository(db)
	outgoingRepository = repository.NewOutgoingRepository(db)
}

func initService() {
	productSvc = services.NewProductService(config, productRepository, incommingRepository, incommingDetailRepository, outgoingRepository)
	incommingSvc = services.NewIncommingService(config, productRepository, incommingRepository, incommingDetailRepository)
	outgoingSvc = services.NewOutgoingService(config, productRepository, outgoingRepository)
}

func runServer() {
	configuration, err := configuration.NewConfiguration()
	if err != nil {
		fmt.Errorf(err.Error())
	}
	config = configuration

	db, err := connectToDatabase(configuration)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	db.AutoMigrate(&dbo.Product{}, &dbo.IncommingProduct{}, &dbo.IncommingProductDetail{}, &dbo.OutgoingProduct{})
	db.Exec("PRAGMA foreign_keys = ON")

	initRepository(db)
	initService()
}

func main() {
	fmt.Println("====================================")
	fmt.Println("Inventory API by Ahmad Reza Musthafa")
	fmt.Println("====================================")

	runServer()
}
