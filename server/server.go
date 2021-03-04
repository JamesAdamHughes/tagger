// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"net/http"
	"tagger/server/handlers"
)

func StartServer() {
	fmt.Printf("\nStarting server at 8080\n")
	http.HandleFunc("/static", handlers.StaticHandler)

	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/callback", handlers.CompleteAuthHandler)
	http.HandleFunc("/playlist", handlers.PlaylistHandler)

	http.HandleFunc("/api/user/playlist", handlers.ApiPlaylistHandler)
	http.HandleFunc("/api/user", handlers.ApiGetUser)

	//http.HandleFunc("/api/track/tag", handlers.ApiSongTagHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
