package client

import (
	"context"
	"golang.org/x/oauth2/clientcredentials"
	"log"
	"net/http"
	"os"
)

//go:generate oapi-codegen -config ../../cfg.yaml ../../cloudapi.yaml

const BASE_URL = "https://fidelis.isams.cloud"

type APIClient struct {
	*http.Client
}

func GetHTTPClient() APIClient {
	clientSecret := os.Getenv("ISAMS_REST_CLIENT_SECRET")

	if clientSecret == "" {
		log.Fatalln("Couldn't find ISAMS_REST_CLIENT_SECRET env var")
	}

	ctx := context.Background()
	conf := clientcredentials.Config{
		ClientID:       "982A09D3-E738-4597-8270-310593223716",
		ClientSecret:   clientSecret,
		TokenURL:       BASE_URL + "/auth/connect/token",
		Scopes:         []string{"restapi"},
		EndpointParams: nil,
		AuthStyle:      0,
	}

	c := conf.Client(ctx)

	return APIClient{c}
}
