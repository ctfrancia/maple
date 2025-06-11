package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ctfrancia/maple/internal/adapters/rest"
	"github.com/ctfrancia/maple/internal/adapters/systemhealth"
	"github.com/ctfrancia/maple/internal/core/services"
)

func main() {
	// Adapters
	hca := systemhealth.NewSystemHealthAdapter()

	// services
	shs := services.NewSystemHealthServicer(hca)

	// Create a new router
	router := rest.NewRouter(shs)

	// Start the server

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println("error starting server: ", err)
		os.Exit(1)
	}
}
