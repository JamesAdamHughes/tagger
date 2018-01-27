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

	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/callback", handlers.CompleteAuthHandler)

	http.HandleFunc("/api/user/playlists", handlers.ApiPlaylistHandlers)
	http.HandleFunc("/api/user", handlers.ApiGetUser)

	http.ListenAndServe(":8080", nil)
}