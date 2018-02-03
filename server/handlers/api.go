package handlers

import (
	"net/http"
	"encoding/json"
	"github.com/zmb3/spotify"
	"fmt"
	"tagger/server/cookies"
	"tagger/server/tags"
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

func ApiSongTagHandler(w http.ResponseWriter, r *http.Request){
	client := cookies.GetClientFromCookies(r)
	
	if r.Method == "POST" {
		var songTag tags.AddSongTagRequest

		// Marshal request into struct
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&songTag)
		if err != nil {
			returnJson(w, ErrorResponse{OK: false, Message: fmt.Sprintf("Something went wrong in post %s", err.Error())})
			return
		}

		// Add user id if cookie set
		songTag.UserId = "0"
		if client != nil {
			user, _ := client.CurrentUser()
			songTag.UserId = user.ID
		}

		// Add the song tag
		err = tags.AddSongTag(&songTag)

		if err != nil {
			returnJson(w, ErrorResponse{OK: false, Message: fmt.Sprintf("Error adding tag %s", err.Error())})
			return
		}

		returnJson(w, ErrorResponse{OK: true, Message: fmt.Sprintf("Success adding tag %+v", songTag)})
	} else if r.Method == "GET" {
		//songId, ok := r.URL.Query()["songId"]; if ok != true {
		//	returnJson(w, ErrorResponse{OK: false, Message: fmt.Sprintf("Missing Param: songId")})
		//	return
		//}
	}

	return
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
