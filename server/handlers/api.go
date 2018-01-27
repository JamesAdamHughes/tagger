package handlers

import (
	"net/http"
	"encoding/json"
	"github.com/zmb3/spotify"
	"fmt"
	"tagger/server/cookies"
)

type PlaylistsResponse struct {
	OK        bool
	Playlists []spotify.SimplePlaylist
}

type PlaylistResponse struct {
	OK bool
	Playlist spotify.FullPlaylist
}

type UserResponse struct {
	OK   bool
	User *spotify.PrivateUser
}

type ErrorResponse struct {
	OK      bool
	Message string
}


func ApiPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	client := cookies.GetClientFromCookies(r)
	if client == nil {
		returnJson(w, ErrorResponse{OK: false, Message: "Error getting client"})
	}

	// If an ID is not provided, return all user playlists
	// Otherwise we get get the specific playlist and all songs
	playlistId, ok := r.URL.Query()["id"]; if ok != true {
		playlistPage, err := client.CurrentUsersPlaylists()
		if err != nil {
			returnJson(w, ErrorResponse{OK: false, Message: err.Error()})
		}

		returnJson(w, PlaylistsResponse{
			OK:        true,
			Playlists: playlistPage.Playlists,
		})

	} else {
		fmt.Println(playlistId)

		user, _ := client.CurrentUser()

		playlist, err := client.GetPlaylist(user.ID, spotify.ID(playlistId[0]))

		currentOffset := 100
		limit := 100

		// Make sure we get all songs, as is a limit of 100tracks per response
		// Loop until all songs retrieved
		for len(playlist.Tracks.Tracks) < playlist.Tracks.Total {

			options := &spotify.Options{
				Limit: &limit,
				Offset: &currentOffset,
			}

			restOfPlaylist, err := client.GetPlaylistTracksOpt(user.ID, spotify.ID(playlistId[0]), options, "")
			if err != nil {
				fmt.Printf("\n\n error %+v", err)
			}

			// Add to the end of the tracks array
			playlist.Tracks.Tracks = append(playlist.Tracks.Tracks, restOfPlaylist.Tracks...)

			currentOffset += 100
		}

		if err != nil || playlist == nil {
			returnJson(w, ErrorResponse{OK: false, Message: fmt.Sprintf("Error getting playlist: %v", playlist)})
			return
		}

		returnJson(w, PlaylistResponse{
			OK:       true,
			Playlist: *playlist,
		})
	}
}


func ApiGetUser(w http.ResponseWriter, r *http.Request) {
	client := cookies.GetClientFromCookies(r)
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
