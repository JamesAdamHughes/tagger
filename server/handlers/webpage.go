package handlers

import (
	"fmt"
	"net/http"
	"tagger/server/spotify"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	client := GetClientFromCookies(r)
	if client != nil {
		//user, err := client.CurrentUser()
		//if err != nil{
		//	fmt.Fprintf(w, "Can't find due to %s", err.Error())
		//}

		//playlists, _ := client.CurrentUsersPlaylists()
		//for _, playlist := range playlists.Playlists {
		//	fmt.Printf("\n%+v\n", playlist.Name)
		//}
		//fmt.Printf("in loaded client")
		page, _ := LoadPage("index_authed")
		fmt.Fprint(w, page.Body)

		//RenderTemplate(w, "index_authed", user)
	} else {
		RenderTemplate(w, "index_anon", Page{})
	}

}

func CompleteAuthHandler(w http.ResponseWriter, r *http.Request){
	client, err := spotify_manager.CompleteAuth(w, r)

	if err != nil {
		RenderTemplate(w, "Error in auth", Page{})
	}

	authToken, _ := client.Token()
	http.SetCookie(w, &http.Cookie{
		Name: "authTokenTMP",
		Value: authToken.AccessToken,
		Expires: authToken.Expiry,
	})

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
