package utils

// import (
// 	"log"
// 	"os"

// 	"golang.org/x/oauth2"
// 	"github.com/coreos/go-oidc"
// 	"github.com/Nerzal/gocloak"
// )

// // InitAuth Initialises the OIDC Client
// func InitAuth() {
//  authEndpoint := os.Getenv("auth_endpoint")
// 	clientID := os.Getenv("auth_client_id")
// 	clientSecret := os.Getenv("auth_client_secret")

// 	ctx := context.Background()
// 	provider, err := oidc.NewProvider(ctx, authEndpoint)
// 	if err != nil {
// 		log.Printf("Couldn't connect to OIDC Provider: %v ", err)
// 	}
// 	// Configure an OpenID Connect aware OAuth2 client.
// 	oauth2Config := oauth2.Config{
// 		ClientID:     clientID,
// 		ClientSecret: clientSecret,
// 		RedirectURL:  "http://localhost:8081/demo/callback",

// 		// Discovery returns the OAuth2 endpoints.
// 		Endpoint: provider.Endpoint(),

// 		// "openid" is a rquired scope for OpenID Connect flows.
// 		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
// 	}

// 	oidcConfig := &oidc.Config{
// 		ClientID: clientID,
// 	}
// 	verifier := provider.Verifier(oidcConfig)
// 	oauth2Config.AuthCodeURL("test")
// }

// var Auth *gocloak.JWT

// func InitAuthEcho() {
// 	authEndpoint := os.Getenv("auth_endpoint")
// 	clientID := os.Getenv("auth_client_id")
// 	clientSecret := os.Getenv("auth_client_secret")
// 	client := gocloak.NewClient(authEndpoint)
// 	Auth, err := client.LoginClient(clientID, clientSecret, "demo")
// 	if err != nil {
// 		log.Panicf("Couldn't generate token: %v", err)
// 	}
// }
