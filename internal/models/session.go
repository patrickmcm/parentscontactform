package models

import (
	"golang.org/x/oauth2"
)

type SessionData struct {
	ISAMSSessionId string
	UserCode       string
	ISAMSToken     ISAMSToken
}

type ISAMSToken struct {
	UserClaims *ClaimExtract
	Token      *oauth2.Token
}
