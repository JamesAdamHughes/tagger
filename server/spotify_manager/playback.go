package spotify_manager

import (
	"fmt"

	"github.com/zmb3/spotify"
)

func GetPlaybackInfo(client *spotify.Client) {
	state, err := client.PlayerState()
	if err != nil {
		fmt.Printf("\n%+v\n", err)
		return
	}
	fmt.Printf("\n%+v\n", state)
}

func PlayerQueueTrack(client *spotify.Client, songId string) error {
	return client.QueueSong(spotify.ID(songId))
}

func PlayerResume(client *spotify.Client) error {
	return client.Play()
}

func PlayerPause(client *spotify.Client) error {
	return client.Pause()
}
