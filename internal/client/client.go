package client

import (
	"net/http"
)

//go:generate oapi-codegen -config ../../cfg.yaml ../../cloudapi.yaml

const BASE_URL = "https://fidelis.isams.cloud"

type APIClient struct {
	*http.Client
}
