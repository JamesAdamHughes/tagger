package cookies

import (
	"net/http"
	"tagger/spotify_manager"
	"golang.org/x/oauth2"
	"github.com/zmb3/spotify"
)

func GetClientFromCookies(r *http.Request) *spotify.Client {
	c, _ := r.Cookie("authTokenTMP")
	if c != nil {
		// get user if possible
		client := spotify_manager.GetClient(&oauth2.Token{
			AccessToken: c.Value,
		})
		return client
	}
	return nil
}
