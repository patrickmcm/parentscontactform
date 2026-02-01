package main

import (
	"embed"
	"fmt"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"io/fs"
	"log"
	"net/http"
	"os"
	"parentscontactform/internal/auth"
	"parentscontactform/internal/client"
	"parentscontactform/internal/handlers"
	"time"
)

//go:embed templates/*.gohtml
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

func main() {
	prod := os.Getenv("PROD")
	if prod == "" {
		setupEnv()
	}

	err := auth.SetupOIDC()
	if err != nil {
		log.Fatalf("Failed to setup OIDC: %v", err)
	}

	if err = sentry.Init(sentry.ClientOptions{
		Dsn: "https://ea3d5a69736ab2a9f530ab0d28a4a9fa@o4510160253747200.ingest.de.sentry.io/4510160288219216",
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for tracing.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
		// Enable structured logs to Sentry
		EnableLogs:    true,
		EnableTracing: true,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	sentryHandler := sentryhttp.New(sentryhttp.Options{})

	hc := client.GetHTTPClient()
	c, err := client.NewClientWithResponses(client.BASE_URL, client.WithHTTPClient(&hc))
	if err != nil {
		log.Fatal("Failed to init iSAMS API Client")
	}

	handler := handlers.NewHandler(staticFS, templateFS, c)

	mux := http.NewServeMux()

	staticSubFS, _ := fs.Sub(staticFS, "static")
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.FS(staticSubFS))))

	mux.HandleFunc("/", sentryHandler.HandleFunc(handler.HandleLoginGet))
	mux.HandleFunc("/login", sentryHandler.HandleFunc(handler.HandleLoginPost))
	mux.HandleFunc("/logout", sentryHandler.HandleFunc(handler.HandleLogoutGet))
	mux.HandleFunc("/form", sentryHandler.HandleFunc(handler.HandleFormGet))
	mux.HandleFunc("/children", sentryHandler.HandleFunc(handler.HandleChildFormGet))
	mux.HandleFunc("/updateChildren", sentryHandler.HandleFunc(handler.HandleChildFormPost))
	mux.HandleFunc("/callback", sentryHandler.HandleFunc(handler.HandleCallback))
	mux.HandleFunc("/submit", sentryHandler.HandleFunc(handler.HandleFormPost))

	srv := http.Server{
		Addr:         os.Getenv("PORT"),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
