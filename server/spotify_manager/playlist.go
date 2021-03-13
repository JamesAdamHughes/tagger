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
	playlist, err := client.GetPlaylist(spotify.ID(playlistID))
	if err != nil {
		return nil, err
	}

	currentOffset := 100
	limit := 100

	// check if the whole playlist response is in redis first, can just return that
	var keyname = fmt.Sprintf("playlist_song_tags_key_%s", playlistID)
	// r := &PlaylistResponse{}
	// err = redis.Get(keyname, r)
	// if err != nil {
	// 	return nil, err
	// }

	// if r.OK {
	// 	return r, nil
	// }

	// Make sure we get all songs, as is a limit of 100tracks per response
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

	playlistSongTags := make(map[string]SongTagsResponse)
	// scrobblerTagger := categoriser.ScrobblerTagger{}
	storedTagger := categoriser.StoredTagger{}

	// Get tags for all songs (should be in the DB if not cached)
	// check if we have stored tags first, otherwise use the scrobbler API
	// for _, track := range playlist.Tracks.Tracks {
	// 	songID := track.Track.ID.String()
	// 	song := categoriser.Song{ID: songID}

	// 	tags, err := storedTagger.GetSongTags(song, user.ID)
	// 	if err == nil && len(tags) > 0 {
	// 		fmt.Printf("\nDATABASE Song: %s, tags: %+v", track.Track.Name, tags)
	// 		playlistSongTags[songID] = SongTagsResponse{
	// 			SongId: songID,
	// 			UserId: user.ID,
	// 			Tags:   tags,
	// 		}
	// 		continue
	// 	} else {
	// 		// pull tags from an API then save them to the DB
	// 		songTags, _ := scrobblerTagger.GetSongTags(categoriser.Song{
	// 			Name:   track.Track.Name,
	// 			Artist: track.Track.Artists[0].Name,
	// 		}, "0")

	// 		fmt.Printf("\n SCROBBLER Song: %s, tags: %+v", track.Track.Name, songTags)

	// 		if len(songTags) > 0 {
	// 			playlistSongTags[songID] = SongTagsResponse{
	// 				SongId: songID,
	// 				UserId: user.ID,
	// 				Tags:   tags,
	// 			}

	// 			// for _, t := range songTags {
	// 			// 	storedTagger.SaveSongTags(song, user.ID, t.Name)
	// 			// }
	// 		}

	// 	}
	// }

	// get all the song data from API concurrently for massive speedup
	tagChan := make(chan SongTagsResponse)
	for _, track := range playlist.Tracks.Tracks {
		song := categoriser.Song{
			Name:   track.Track.Name,
			Artist: track.Track.Artists[0].Name,
			ID:     track.Track.ID.String(),
		}
		go getSongTagsWorker(tagChan, song)
	}

	for j := 1; j <= len(playlist.Tracks.Tracks); j++ {
		songTagResponse := <-tagChan
		playlistSongTags[songTagResponse.SongId] = songTagResponse
		for _, t := range songTagResponse.Tags {
			storedTagger.SaveSongTags(songTagResponse.SongId, user.ID, t.Name)
		}
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

func getSongTagsWorker(tagsChannel chan SongTagsResponse, song categoriser.Song) {
	scrobblerTagger := categoriser.ScrobblerTagger{}
	// pull tags from an API then save them to the DB
	songTags, _ := scrobblerTagger.GetSongTags(song, "0")
	fmt.Printf("in worker for song %s %s\n", song.Name, song.ID)
	var songTagResponse = SongTagsResponse{
		SongId: song.ID,
		Tags:   songTags,
	}

	tagsChannel <- songTagResponse
}
