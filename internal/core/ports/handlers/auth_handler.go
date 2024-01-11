package handlers

import "net/http"

type AuthHandlers interface {
	SignInHandler(w http.ResponseWriter, req *http.Request)
	SignUpHandler(w http.ResponseWriter, req *http.Request)
}
