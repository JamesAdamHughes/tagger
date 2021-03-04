package handlers

import (
	"fmt"
	"net/http"
)

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("going static")
}