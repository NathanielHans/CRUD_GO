package main

import (
	"CRUD-go/config"
	categoryhandlers "CRUD-go/handlers/categoryHandlers"
	producthandlers "CRUD-go/handlers/productHandlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	config.ConnectDB()

	r := mux.NewRouter()

	// Path-based API
	r.HandleFunc("/categories", categoryhandlers.CategoryIndex).Methods("GET")
	r.HandleFunc("/categories", categoryhandlers.CategoryStore).Methods("POST")
	r.HandleFunc("/categories/{id}", categoryhandlers.CategoryFindByID).Methods("GET")
	r.HandleFunc("/categories/{id}", categoryhandlers.CategoryUpdate).Methods("PUT")
	r.HandleFunc("/categories/{id}", categoryhandlers.CategoryDelete).Methods("DELETE")

	r.HandleFunc("/products", producthandlers.ProductIndex).Methods("GET")
	r.HandleFunc("/products/{id}", producthandlers.ProductFindByID).Methods("GET")
	r.HandleFunc("/products", producthandlers.ProductStore).Methods("POST")
	r.HandleFunc("/products/{id}", producthandlers.ProductUpdate).Methods("PUT")
	r.HandleFunc("/products/{id}", producthandlers.ProductDelete).Methods("DELETE")
	log.Println("Server Running on port :8081")
	http.ListenAndServe(":8081", r)
}
