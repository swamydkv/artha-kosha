package main

import (
	"log"
	"net/http"
	"time"

	"artha-kosha/apps/finance-api/internal/auth"
	"artha-kosha/apps/finance-api/internal/config"
	internalhttp "artha-kosha/apps/finance-api/internal/http"
)

func main() {
	cfg := config.Load()
	provider := auth.NewLocalAuthProvider()
	mux := internalhttp.NewRouter(provider)
	// wrap with middleware: logging + timeout

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("starting auth API on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
