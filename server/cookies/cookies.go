package cookies

import (
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"net/http"
	"tagger/server/spotify_manager"
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
