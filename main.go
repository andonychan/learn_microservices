package main

import (
	"context"
	"github.com/gorilla/mux"
	"intro/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	ph := handlers.NewProducts(l)

	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	putRouter := sm.Methods(http.MethodPut).Subrouter()
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	//deleteRouter := sm.Methods("DELETE").Subrouter()

	getRouter.HandleFunc("/", ph.GetProducts)
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)
	postRouter.HandleFunc("/", ph.AddProduct)

	putRouter.Use(ph.MiddlewareProductValidation, ph.MiddlewareProductIDValidation)
	postRouter.Use(ph.MiddlewareProductValidation)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  2 * time.Minute,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt)
	signal.Notify(sigchan, os.Kill)

	sig := <-sigchan
	l.Println("shutting down", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := s.Shutdown(tc)
	if err != nil {
		l.Fatal(err)
	}

}
