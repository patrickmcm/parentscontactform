package auth

import (
	"context"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"log"
	"os"
)

var (
	Cfg      *oauth2.Config
	Verifier *oidc.IDTokenVerifier
)

func SetupOIDC() error {
	ctx := context.Background()
	domain := "https://fidelis.isams.cloud"
	provider, err := oidc.NewProvider(ctx, domain+"/auth")
	if err != nil {
		return err
	}

	redirect := os.Getenv("REDIRECT_URI")

	clientSecret := os.Getenv("ISAMS_OIDC_CLIENT_SECRET")

	if clientSecret == "" {
		log.Fatalln("Couldn't find ISAMS_OIDC_CLIENT_SECRET env var")
	}

	Cfg = &oauth2.Config{
		ClientID:     "Fidelis.Parents.UpdateForm",
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  redirect,
		Scopes:       []string{"openid", "profile", "email", "iSAMS.CloudPortals.Api"},
	}

	Verifier = provider.Verifier(&oidc.Config{
		ClientID: "Fidelis.Parents.UpdateForm",
	})
	return nil
}
