// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"net/http"
	"tagger/server/handlers"
)

const PORT = 8081

func StartServer() {
	fmt.Printf("\nStarting server from docker %d\n", PORT)
	http.HandleFunc("/static", handlers.StaticHandler)

	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/callback", handlers.CompleteAuthHandler)
	http.HandleFunc("/playlist", handlers.PlaylistHandler)

	http.HandleFunc("/api/user/playlist", handlers.ApiPlaylistHandler)
	http.HandleFunc("/api/user", handlers.ApiGetUser)
	http.HandleFunc("/api/player", handlers.ApiGetPlayer)
	http.HandleFunc("/api/player/queuetrack", handlers.ApiPlayerQueueTrack)

	//http.HandleFunc("/api/track/tag", handlers.ApiSongTagHandler)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil); err != nil {
		panic(err)
	}
}
