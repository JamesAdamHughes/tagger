package categoriser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type lastFmTagResponse struct {
	TopTags `json:"toptags"`
}

type TopTags struct {
	Tag []LastFmTag `json:"tag"`
}

type LastFmTag struct {
	Count int64
	Name  string
}

type Tag struct {
	ID   int64
	Name string
}

type Song struct {
	Name   string
	Artist string
	ID     string
}

type Tagger interface {
	GetSongTags(song Song, userId string) (tags []Tag, err error)
	SaveSongTags(song Song, userId string, tagName string) (err error)
}

type ScrobblerTagger struct{}

func (gt ScrobblerTagger) GetSongTags(song Song, userId string) (tags []Tag, err error) {

	song.Artist = strings.Replace(song.Artist, " ", "+", -1)
	song.Name = strings.Replace(song.Name, " ", "+", -1)

	resp, err := http.Get(fmt.Sprintf("http://ws.audioscrobbler.com/2.0/?method=track.gettoptags&artist=%s&track=%s&api_key=0e367bb3137650849d8cc15446ed3566&format=json", song.Artist, song.Name))
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var lastfmResponse lastFmTagResponse
	err = json.Unmarshal(body, &lastfmResponse)
	if err != nil {
		panic(err)
	}

	log.Println(lastfmResponse)

	for idx, tag := range lastfmResponse.TopTags.Tag {
		if idx > 4 {
			continue
		}
		tags = append(tags, Tag{
			ID:   0,
			Name: tag.Name,
		})
	}

	return tags, nil
}

func (gt ScrobblerTagger) SaveSongTags(song Song, userId string, tagName string) (err error) {
	return nil
}

type StoredTagger struct{}

func (st StoredTagger) GetSongTags(song Song, userId string) (tags []Tag, err error) {
	r := GetSongTagRequest{song.ID, userId}
	tags, err = GetSongTags(r)

	if err != nil {
		fmt.Printf("\nStoredTagger error %s, %+v\n", err.Error(), tags)

		return nil, err
	}

	fmt.Printf("\nStoredTagger %s, %+v\n", song.ID, tags)

	return tags, nil
}

func (st StoredTagger) SaveSongTags(song Song, userId string, tagName string) (err error) {
	r := AddSongTagRequest{SongId: song.ID, UserId: userId, TagName: tagName}
	err = AddPlaylistSongTag(r)
	if err != nil {
		return err
	}
	return nil
}
