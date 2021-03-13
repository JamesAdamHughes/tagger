package spotify_manager

import (
	"fmt"
	"tagger/server/categoriser"
	"tagger/server/redis"

	"github.com/zmb3/spotify"
)

type PlaylistResponse struct {
	OK           bool
	Playlist     spotify.FullPlaylist
	PlaylistTags map[string]SongTagsResponse
}

type SongTagsResponse struct {
	SongId string
	UserId string
	Tags   []categoriser.Tag
}

func FetchSongsFromPlaylist(client *spotify.Client, playlistID string) (playlistResponse *PlaylistResponse, err error) {
	user, _ := client.CurrentUser()

	// check if the whole playlist response is in redis first, can just return that
	var keyname = fmt.Sprintf("playlist_song_tags_key_%s", playlistID)
	r := &PlaylistResponse{}
	err = redis.Get(keyname, r)
	if err != nil {
		return nil, err
	}

	if r.OK {
		fmt.Printf("\nRedis won lads\n")
		// return r, nil
	}

	playlist, err := getSpotifyPlaylist(client, playlistID)
	if err != nil {
		return nil, err
	}

	playlistSongTags := make(map[string]SongTagsResponse)
	tagChan := make(chan SongTagsResponse)
	var jobCounter int

	// Get tags for all songs using go routines for speed
	for _, track := range playlist.Tracks.Tracks {
		songID := track.Track.ID.String()
		song := categoriser.Song{
			Name:   track.Track.Name,
			Artist: track.Track.Artists[0].Name,
			ID:     songID,
		}

		go getSavedSongTagsWorker(tagChan, song, user.ID)
		jobCounter++
	}

	for j := 1; j <= jobCounter; j++ {
		songTagResponse := <-tagChan
		playlistSongTags[songTagResponse.SongId] = songTagResponse
		fmt.Printf("done with song %s\n", songTagResponse.SongId)
	}
	close(tagChan)

	playlistResponse = &PlaylistResponse{
		OK:           true,
		Playlist:     *playlist,
		PlaylistTags: playlistSongTags,
	}

	// cache result in redis
	err = redis.Set(keyname, &playlistResponse, 60*60)
	if err != nil {
		return nil, err
	}

	return playlistResponse, nil
}

// getSavedSongTagsWorker Get tags for all songs (should be in the DB if not cached)
// check if we have stored tags first, otherwise use the scrobbler API
func getSavedSongTagsWorker(tagsChannel chan SongTagsResponse, song categoriser.Song, userID string) {
	scrobblerTagger := categoriser.ScrobblerTagger{}
	storedTagger := categoriser.StoredTagger{}

	// first check if in the db, then go to scribbler
	songTags, err := storedTagger.GetSongTags(song, userID)
	inDB := true
	if err == nil && len(songTags) == 0 {
		fmt.Printf("in  saved  worker was NOT IN DB for song %s %s\n", song.Name, song.ID)
		inDB = false
		songTags, _ = scrobblerTagger.GetSongTags(song, "0")
	}

	// pull tags from database then save them to the DB
	fmt.Printf("in  saved  worker for song %s %s\n", song.Name, song.ID)
	var songTagResponse = SongTagsResponse{
		SongId: song.ID,
		Tags:   songTags,
	}

	tagsChannel <- songTagResponse

	// run the DB operation after sending the response back, as we don't need to block on this. Rather return to user faster
	if !inDB {
		for _, t := range songTags {
			storedTagger.SaveSongTags(songTagResponse.SongId, userID, t.Name)
		}
	}
}

func getSpotifyPlaylist(client *spotify.Client, playlistID string) (*spotify.FullPlaylist, error) {
	playlist, err := client.GetPlaylist(spotify.ID(playlistID))
	if err != nil {
		return nil, err
	}

	currentOffset := 100
	limit := 100

	// Make sure we get all songs from spotify, as is a limit of 100tracks per response
	// Loop until all songs retrieved
	for len(playlist.Tracks.Tracks) < playlist.Tracks.Total {

		options := &spotify.Options{
			Limit:  &limit,
			Offset: &currentOffset,
		}

		restOfPlaylist, err := client.GetPlaylistTracksOpt(spotify.ID(playlistID), options, "")
		if err != nil {
			fmt.Printf("\n\n error %+v", err)
		}

		// Add to the end of the tracks array
		playlist.Tracks.Tracks = append(playlist.Tracks.Tracks, restOfPlaylist.Tracks...)

		currentOffset += 100
	}

	return playlist, nil
}
