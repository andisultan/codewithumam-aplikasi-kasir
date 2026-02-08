package main

import (
	"aplikasi-kasir/database"
	"aplikasi-kasir/handlers"
	"aplikasi-kasir/repositories"
	"aplikasi-kasir/services"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database", err)
	}

	defer db.Close()

	// Product
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	http.HandleFunc("/api/products", productHandler.HandleProducts)
	http.HandleFunc("/api/products/", productHandler.HandleProductByID)

	// Product Category
	productCategoryRepo := repositories.NewProductCategoryRepository(db)
	productCategoryService := services.NewProductCategoryService(productCategoryRepo)
	productCategoryHandler := handlers.NewProductCategoryHandler(productCategoryService)

	http.HandleFunc("/api/product-categories", productCategoryHandler.HandleProductCategories)
	http.HandleFunc("/api/product-categories/", productCategoryHandler.HandleProductCategoryByID)

	// Transaction
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	http.HandleFunc("/api/checkout", transactionHandler.HandleCheckout) // POST

	// Report
	reportRepo := repositories.NewReportRepository(db)
	reportService := services.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)

	http.HandleFunc("/api/report/hari-ini", reportHandler.HandleDailyReport)
	http.HandleFunc("/api/report", reportHandler.HandleReport)

	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running on port", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
