package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"intro/product-api/data"
	"log"
	"net/http"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")
	lp := data.GetProduct()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Error marshalling product", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")

	//Chain from middleware
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	p.l.Printf("Prod: %#v", prod)
	data.AddProduct(prod)
}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//id, err := strconv.Atoi(vars["id"])
	//if err != nil {
	//	http.Error(rw, "Error converting product ID to int", http.StatusBadRequest)
	//	return
	//}

	id := r.Context().Value(KeyProductID{}).(int)

	p.l.Println("Handle PUT Product")

	//Chain from middleware
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	err := data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Error updating product", http.StatusInternalServerError)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
	}
}

type KeyProduct struct{}

func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}
		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Error unmarshalling product", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

type KeyProductID struct{}

func (p *Products) MiddlewareProductIDValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(rw, "Error converting product ID to int", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProductID{}, id)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
