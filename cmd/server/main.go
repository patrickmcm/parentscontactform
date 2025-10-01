package main

import (
	"log"
	"net/http"
	"os"
	"parentscontactform/internal/auth"
	"parentscontactform/internal/client"
	"parentscontactform/internal/handlers"
)

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

	handlers.RestClient = c

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("/Users/patrickmcm/GolandProjects/parentscontactform/cmd/server/static"))))

	http.HandleFunc("/", handlers.HandleLoginGet)
	http.HandleFunc("/login", handlers.HandleLoginPost)
	http.HandleFunc("/logout", handlers.HandleLogoutGet)
	http.HandleFunc("/form", handlers.HandleFormGet)
	http.HandleFunc("/children", handlers.HandleChildFormGet)
	http.HandleFunc("/updateChildren", handlers.HandleChildFormPost)
	http.HandleFunc("/callback", handlers.HandleCallback)
	http.HandleFunc("/submit", handlers.HandleFormPost)

	http.ListenAndServe(":3000", nil)
}
