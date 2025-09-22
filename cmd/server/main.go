package main

import (
	"log"
	"net/http"
	"os"
	"parentscontactform/internal/auth"
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

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("C:\\Users\\Patrick\\GolandProjects\\parentscontactform\\cmd\\server\\static"))))

	http.HandleFunc("/", handlers.HandleLoginGet)
	http.HandleFunc("/login", handlers.HandleLoginPost)
	http.HandleFunc("/logout", handlers.HandleLogoutGet)
	http.HandleFunc("/form", handlers.HandleFormGet)
	http.HandleFunc("/callback", handlers.HandleCallback)
	http.HandleFunc("/submit", handlers.HandleFormPost)

	http.ListenAndServe(":3000", nil)
}
