package main

import (
	"context"
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

	sm := http.NewServeMux()

	sm.Handle("/", ph)

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
