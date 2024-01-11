package handlers

import (
	"net/http"
)

type BondHandlers interface {
	ListBondsHandler(w http.ResponseWriter, req *http.Request)
	GetBondByUUIDHandler(w http.ResponseWriter, req *http.Request)
	CreateBondHandler(w http.ResponseWriter, req *http.Request)
	UpdateBondHandler(w http.ResponseWriter, req *http.Request)
	DeleteBondHandler(w http.ResponseWriter, req *http.Request)
}
