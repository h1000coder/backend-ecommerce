package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"soulstreet/internal/config/cors"
	"soulstreet/internal/config/db"
	"soulstreet/internal/handler"
	"soulstreet/internal/repository"
	"soulstreet/internal/service"
)


func init() {
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		err := os.MkdirAll("uploads", 0755)
		if err != nil {
			log.Fatalf("Error creating uploads directory: %v", err)
		}
	}
}

func main() {
	
	dbConn := db.ConnectDB()
	err := db.CreateTableIfNotExist(dbConn)
	if err != nil {
		log.Fatal(err)
	}

	// Inicializa o repositório, serviço e manipulador para produtos
	repoProduct := repository.NewProductRepositoryDB(dbConn)
	serviceProduct := service.NewProductService(repoProduct)
	handlerProduct := handler.NewProductHandler(*serviceProduct)
	
	// Inicializa o repositório, serviço e manipulador para pagamentos
	repoPayment := repository.NewPaymentRepositoryDB(dbConn)
	servicePayment := service.NewPaymentService(repoPayment)
	handlerPayment := handler.NewPaymentHandler(*servicePayment)


	mux := http.NewServeMux()
	
	// Endpoints produto
	mux.HandleFunc("POST /product", handlerProduct.CreateProduct)
	mux.HandleFunc("GET /products", handlerProduct.GetAll)
	mux.HandleFunc("GET /product", handlerProduct.GetByID)
	mux.Handle("GET /images/", http.StripPrefix("/images/", http.FileServer(http.Dir("uploads"))))
	mux.HandleFunc("DELETE /product", handlerProduct.DeleteProduct)
	mux.HandleFunc("GET /product/name", handlerProduct.GetProductByName)

	// Endpoints pagamento
	mux.HandleFunc("POST /webhook", handlerPayment.CreatePayment)

	
	fmt.Println("Servidor rodando http://localhost:8080")
	
	handlerCors := cors.CORS()(mux)
	
	http.ListenAndServe(":8080", handlerCors)
}