package client

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"parentscontactform/internal/models"
)

func GetContactInfo(sessionId string) ([]models.ParentContactInfo, error) {
	var baseUrl string

	if baseUrl = os.Getenv("BASE-URL"); baseUrl == "" {
		baseUrl = "https://fidelis.isams.cloud/Legacy/Api/Rest/1.0"
	}

	endpoint := baseUrl + "/portals/parent/comms/contacts"

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		var temp []models.ParentContactInfo
		return temp, err
	}

	req.Header.Add("sessionId", sessionId)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		var temp []models.ParentContactInfo
		return temp, err
	}

	responseBodyBuf, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		var temp []models.ParentContactInfo
		return temp, err
	}

	var responseBody []models.ParentContactInfo

	err = json.Unmarshal(responseBodyBuf, &responseBody)
	if err != nil {
		var temp []models.ParentContactInfo
		return temp, err
	}

	return responseBody, nil
}

func GetUserAccountDetails(data *models.SessionData) (models.CurrentUserInfo, error) {
	var baseUrl string

	if baseUrl = os.Getenv("BASE-URL"); baseUrl == "" {
		baseUrl = "https://fidelis.isams.cloud/Legacy/Api/Rest/1.0"
	}

	endpoint := baseUrl + "/users/user/" + data.UserCode

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		var temp models.CurrentUserInfo
		return temp, err
	}

	req.Header.Add("sessionId", data.ISAMSSessionId)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		var temp models.CurrentUserInfo
		return temp, err
	}

	responseBodyBuf, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		var temp models.CurrentUserInfo
		return temp, err
	}

	var responseBody models.CurrentUserInfo

	err = json.Unmarshal(responseBodyBuf, &responseBody)
	if err != nil {
		var temp models.CurrentUserInfo
		return temp, err
	}

	return responseBody, nil
}

func GetUserChildren(data *models.SessionData) ([]models.ChildInfo, error) {
	var baseUrl string

	if baseUrl = os.Getenv("BASE-URL"); baseUrl == "" {
		baseUrl = "https://fidelis.isams.cloud/Legacy/Api/Rest/1.0"
	}

	endpoint := baseUrl + "/portals/parent/children"

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		var temp []models.ChildInfo
		return temp, err
	}

	req.Header.Add("sessionId", data.ISAMSSessionId)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		var temp []models.ChildInfo
		return temp, err
	}

	responseBodyBuf, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		var temp []models.ChildInfo
		return temp, err
	}

	var responseBody []models.ChildInfo

	err = json.Unmarshal(responseBodyBuf, &responseBody)
	if err != nil {
		var temp []models.ChildInfo
		return temp, err
	}

	return responseBody, nil
}
