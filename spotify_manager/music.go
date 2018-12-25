package spotify_manager

import (
	"fmt"
	"github.com/zmb3/spotify"
	"tagger/categoriser"
)

type PlaylistResponse struct {
	OK bool
	Playlist spotify.FullPlaylist
	PlaylistTags []SongTagsResponse
}

type SongTagsResponse struct {
	SongId string
	UserId string
	Tags  []categoriser.Tag
}

func FetchSongsFromPlaylist(client *spotify.Client, playlistID string) (*PlaylistResponse, error) {
	user, _ := client.CurrentUser()

	playlist, err := client.GetPlaylist(user.ID, spotify.ID(playlistID))
	if err != nil {
		return nil, err
	}

	currentOffset := 100
	limit := 100

	// Make sure we get all songs, as is a limit of 100tracks per response
	// Loop until all songs retrieved
	for len(playlist.Tracks.Tracks) < playlist.Tracks.Total {

		options := &spotify.Options{
			Limit: &limit,
			Offset: &currentOffset,
		}

		restOfPlaylist, err := client.GetPlaylistTracksOpt(user.ID, spotify.ID(playlistID), options, "")
		if err != nil {
			fmt.Printf("\n\n error %+v", err)
		}

		// Add to the end of the tracks array
		playlist.Tracks.Tracks = append(playlist.Tracks.Tracks, restOfPlaylist.Tracks...)

		currentOffset += 100
	}

	playlistSongTags := []SongTagsResponse{}
	genreTagger := categoriser.GenreTagger{}

	// Get tags for all songs (should be cached)
	for _, track := range playlist.Tracks.Tracks {

		songTags, _ := genreTagger.GetSongTags(categoriser.Song{
			Name: track.Track.Name,
			Artist: track.Track.Artists[0].Name,
		}, "0")

		fmt.Printf("\n Song: %s, tags: %+v", track.Track.Name, songTags)

		//tags, err := tags.GetSongTags(tags.GetSongTagRequest{
		//	SongId: track.Track.ID.String(),
		//	UserId: user.ID,
		//})
		//if err != nil {
		//	fmt.Printf("\n\n error %+v", err)
		//}
		if len(songTags) > 0 {
			playlistSongTags = append(playlistSongTags, SongTagsResponse{
				SongId: track.Track.ID.String(),
				UserId: user.ID,
				Tags: songTags,
			})
		}
	}

	return &PlaylistResponse{
		OK:true,
		Playlist: *playlist,
		PlaylistTags:playlistSongTags,
	}, nil
}
