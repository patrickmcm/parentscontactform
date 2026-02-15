package models

type LoginBody struct {
	ApiKey   string `json:"apiKey"`
	Password string `json:"password"`
	Portal   string `json:"portal"`
	UserName string `json:"userName"`
}

type LoginResponse struct {
	SessionId string `json:"sessionId"`
	UserCode  string `json:"userCode"`
}

type ClaimExtract struct {
	SessionId string `json:"sid"`
	UserCode  string `json:"sub"`
	ISAMSGuid string `json:"iSAMS.PersonId"`
	Email     string `json:"email"`
	Forename  string `json:"name"`
	Surname   string `json:"family_name"`
}

type LogoutToken struct {
	SessionId string                 `json:"sid"`
	Events    map[string]interface{} `json:"events"`
}
