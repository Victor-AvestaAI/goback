package main

import (
	"log"
	"net/http"

	"github.com/Victor-AvestaAI/goback/internal/admin"
	"github.com/Victor-AvestaAI/goback/internal/api"
)

const port string = "8080"

func main() {

	mux := http.NewServeMux()

	baseHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))

	apiCfg := api.ApiConfig{}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(baseHandler))

	mux.HandleFunc("GET /api/healthz", admin.HandlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", admin.HandlerValidation)

	mux.HandleFunc("GET /admin/metrics", apiCfg.HandlerPrintRequestCount)
	mux.HandleFunc("POST /admin/reset", apiCfg.HandlerResetRequestCount)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())

}
