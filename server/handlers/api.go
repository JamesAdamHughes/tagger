package handlers

import (
	"fmt"
	"github.com/zmb3/spotify"
	"net/http"
	"tagger/server/cookies"
	"tagger/server/spotify_manager"
)

type PlaylistsResponse struct {
	OK        bool
	Playlists []spotify.SimplePlaylist
}

type PlaylistResponse struct {
	OK           bool
	Playlist     spotify.FullPlaylist
	PlaylistTags []SongTagsResponse
}

type UserResponse struct {
	OK   bool
	User *spotify.PrivateUser
}

type SongTagsResponse struct {
	SongId string
	UserId string
	Tags   []Tag
}

type Tag struct {
	ID   int64
	Name string
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
	playlistId, ok := r.URL.Query()["id"]
	if ok != true {
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

		res, err := spotify_manager.FetchSongsFromPlaylist(client, playlistId[0])
		if err != nil {
			panic(err)
		}

		returnJson(w, res)
	}
}

//func ApiSongTagHandler(w http.ResponseWriter, r *http.Request){
//	client := cookies.GetClientFromCookies(r)
//
//	if r.Method == "POST" {
//		var songTag tags.AddSongTagRequest
//
//		// Marshal request into struct
//		decoder := json.NewDecoder(r.Body)
//		err := decoder.Decode(&songTag)
//		if err != nil {
//			returnError(w, ErrorResponse{OK: false, Message: fmt.Sprintf("Something went wrong in post %s", err.Error())})
//			return
//		}
//
//		// Add user id if cookie set
//		songTag.UserId = "0"
//		if client != nil {
//			user, _ := client.CurrentUser()
//			songTag.UserId = user.ID
//		}
//
//		// Add the song tag
//		err = tags.AddSongTag(&songTag)
//
//		if err != nil {
//			returnJson(w, ErrorResponse{OK: false, Message: fmt.Sprintf("Error adding tag %s", err.Error())})
//			return
//		}
//
//		returnJson(w, songTag)
//	} else if r.Method == "GET" {
//		songIdString, ok := r.URL.Query()["songId"]; if ok != true {
//			returnJson(w, ErrorResponse{OK: false, Message: fmt.Sprintf("Error getting tag. Missing Param: songId")})
//			return
//		}
//
//		request := tags.GetSongTagRequest{
//			SongId: songIdString[0],
//			UserId: "0",
//		}
//
//		if client != nil {
//			user, _ := client.CurrentUser()
//			request.UserId = user.ID
//		}
//
//		songTags, err := tags.GetSongTags(request)
//		if err != nil {
//			returnJson(w, ErrorResponse{OK: false, Message: fmt.Sprintf("Error getting tag %s", err.Error())})
//			return
//		}
//
//		returnJson(w, SongTagsResponse{
//			SongId: request.SongId,
//			UserId: request.UserId,
//			Tags:   songTags,
//		})
//	}
//
//	return
//}

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

func ApiGetPlayer(w http.ResponseWriter, r *http.Request) {
	client := cookies.GetClientFromCookies(r)
	if client == nil {
		returnJson(w, ErrorResponse{OK: false, Message: "Error getting client"})
	}

	user, _ := client.CurrentUser()

	returnJson(w, UserResponse{
		OK:   true,
		User: user,
	})

	spotify_manager.GetPlaybackInfo(client)
}

func ApiPlayerQueueTrack(w http.ResponseWriter, r *http.Request) {
	client := cookies.GetClientFromCookies(r)
	if client == nil {
		returnJson(w, ErrorResponse{OK: false, Message: "Error getting client"})
	}

	songId, ok := r.URL.Query()["id"]

	fmt.Printf("\n%+v\n", songId)

	if !ok {
		returnJson(w, ErrorResponse{OK: false, Message: "Missing song ID"})
	}

	spotify_manager.PlayerQueueTrack(client, songId[0])

	returnJson(w, ErrorResponse{OK: true, Message: "hello"})
}
