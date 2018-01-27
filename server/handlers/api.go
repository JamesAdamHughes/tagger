package handlers

import (
	"net/http"
	"encoding/json"
	"github.com/zmb3/spotify"
	"fmt"
)

type PlaylistResponse struct {
	OK        bool
	Playlists []spotify.SimplePlaylist
}

type UserResponse struct {
	OK   bool
	User *spotify.PrivateUser
}

type ErrorResponse struct {
	OK      bool
	Message string
}

func ApiPlaylistHandlers(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n\nhello this is api\n\n")

	client := GetClientFromCookies(r)
	if client == nil {
		returnJson(w, ErrorResponse{OK: false, Message: "Error getting client"})
	}

	playlistPage, err := client.CurrentUsersPlaylists()
	if err != nil {
		returnJson(w, ErrorResponse{OK: false, Message: err.Error()})
	}

	returnJson(w, PlaylistResponse{
		OK:        true,
		Playlists: playlistPage.Playlists,
	})
}

func ApiGetUser(w http.ResponseWriter, r *http.Request) {
	client := GetClientFromCookies(r)
	if client == nil {
		returnJson(w, ErrorResponse{OK: false, Message: "Error getting client"})
	}

	user, _ := client.CurrentUser()

	returnJson(w, UserResponse{
		OK:   true,
		User: user,
	})

}

func returnJson(w http.ResponseWriter, i interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(i)
	return
}
