package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"parentscontactform/internal/auth"
	"parentscontactform/internal/client"
	"parentscontactform/internal/handlers"
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

	hc := client.GetHTTPClient()
	c, err := client.NewClientWithResponses(client.BASE_URL, client.WithHTTPClient(&hc))
	if err != nil {
		log.Fatal("Failed to init iSAMS API Client")
	}

	handler := handlers.NewHandler(staticFS, templateFS, c)

	staticSubFS, _ := fs.Sub(staticFS, "static")
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.FS(staticSubFS))))

	http.HandleFunc("/", handler.HandleLoginGet)
	http.HandleFunc("/login", handler.HandleLoginPost)
	http.HandleFunc("/logout", handler.HandleLogoutGet)
	http.HandleFunc("/form", handler.HandleFormGet)
	http.HandleFunc("/children", handler.HandleChildFormGet)
	http.HandleFunc("/updateChildren", handler.HandleChildFormPost)
	http.HandleFunc("/callback", handler.HandleCallback)
	http.HandleFunc("/submit", handler.HandleFormPost)

	err = http.ListenAndServe(os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
