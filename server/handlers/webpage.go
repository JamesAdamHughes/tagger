package handlers

import (
	"fmt"
	"net/http"
	"tagger/server/spotify"
	"tagger/server/cookies"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	client := cookies.GetClientFromCookies(r)

	if client != nil {
		page, _ := LoadPage("index_authed")
		fmt.Fprint(w, page.Body)
	} else {
		fmt.Println("User is anonymous")
		RenderTemplate(w, "index_anon", Page{})
	}
}

func PlaylistHandler(w http.ResponseWriter, r *http.Request) {
	client := cookies.GetClientFromCookies(r)
	if client != nil {
		page, _ := LoadPage("playlist")
		fmt.Fprint(w, page.Body)
	} else {
		RenderTemplate(w, "not_authed", Page{})
	}
}

func CompleteAuthHandler(w http.ResponseWriter, r *http.Request){
	client, err := spotify_manager.CompleteAuth(w, r)

	if err != nil {
		RenderTemplate(w, "Error in auth", Page{})
	}

	// Set a cookie in the browser for authorization
	authToken, _ := client.Token()
	http.SetCookie(w, &http.Cookie{
		Name: "authTokenTMP",
		Value: authToken.AccessToken,
		Expires: authToken.Expiry,
	})

	// Redirect back to the index page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	authDetails := spotify_manager.GetAuthDetails()

	// Generate auth url with state
	url := authDetails.Auth.AuthURL(authDetails.State)

	// Show link to user to auth with spotify
	RenderTemplate(w, "register", AuthPage{
		Title: "Register with Spotify",
		Body: "Rgeister with your Spotify account to use tagger",
		AuthUrl: url,
	})
}
