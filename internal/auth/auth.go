package auth

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"parentscontactform/internal/models"
)

func LoginUser(email string, password string) (models.SessionData, error) {
	var baseUrl string

	if baseUrl = os.Getenv("BASE-URL"); baseUrl == "" {
		baseUrl = "https://fidelis.isams.cloud/Legacy/Api/Rest/1.0"
	}

	endpoint := baseUrl + "/portals/New%20Portal%20API%20Key/login"

	body := models.LoginBody{
		ApiKey:   os.Getenv("API-KEY"),
		Password: password,
		Portal:   "New Portal API Key",
		UserName: email,
	}

	encodedBody, err := json.Marshal(body)
	if err != nil {
		var temp models.SessionData
		return temp, err
	}

	encodedBodyBuffer := bytes.NewBuffer(encodedBody)

	resp, err := http.Post(endpoint, "application/json", encodedBodyBuffer)
	if err != nil {
		var temp models.SessionData
		return temp, err
	}

	responseBodyBuf, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		var temp models.SessionData
		return temp, err
	}

	var responseBody models.LoginResponse

	err = json.Unmarshal(responseBodyBuf, &responseBody)
	if err != nil {
		var temp models.SessionData
		return temp, err
	}

	session := models.SessionData{
		ISAMSSessionId: responseBody.SessionId,
		UserCode:       responseBody.UserCode,
	}

	return session, nil

}
