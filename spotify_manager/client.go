package spotify_manager

import (
	"github.com/zmb3/spotify"
	"net/http"
	"log"
	"fmt"
	//"os"
	"golang.org/x/oauth2"
	"os"
)

type AuthDetails struct {
	Auth spotify.Authenticator
	State string
}

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectURI = "http://localhost:8080/callback"
//var Client spotify_manager.Client

type Client spotify.Client

func init() {
	os.Setenv("SPOTIFY_ID", "f346f777add648db9f09e7c6ddf87f34")
	os.Setenv("SPOTIFY_SECRET", "1b5bd7bd2e084a178e77e24cf4940f21")
}

func GetAuthDetails() AuthDetails {
	return AuthDetails{
		Auth: spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadPrivate),
		State: "abc123",
	}
}

func GetClient(token *oauth2.Token) *spotify.Client {
	authDetails := GetAuthDetails()
	client := authDetails.Auth.NewClient(token)
	return &client
}

// Handles OAuth2 with Spotify
// Returns the spotify_manager client
func CompleteAuth(w http.ResponseWriter, r *http.Request) (*spotify.Client, error) {
	authDetails := GetAuthDetails()

	tok, err := authDetails.Auth.Token(authDetails.State, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != authDetails.State {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, authDetails.State)
	}
	// use the token to get an authenticated client

	// Todo save this token for repeat calls
	// todo save this in user session
	Client := authDetails.Auth.NewClient(tok)

	// use the client to make calls that require authorization
	user, err := Client.CurrentUser()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	fmt.Println("You are logged in as:", user.ID)

	return &Client, nil
}



