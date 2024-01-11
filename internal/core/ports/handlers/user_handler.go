package handlers

import "net/http"

type UserHandlers interface {
	GetProfileHandler(w http.ResponseWriter, req *http.Request)
	UpdateProfileHandler(w http.ResponseWriter, req *http.Request)
	GetUserBondsHandler(w http.ResponseWriter, req *http.Request)
}
