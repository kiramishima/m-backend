package handlers

import (
	"net/http"
)

// MarketBondsHandlers interface
type MarketBondsHandlers interface {
	ListMarketBondsHandler(w http.ResponseWriter, req *http.Request)
	GetMarketBondByIDHandler(w http.ResponseWriter, req *http.Request)
	BuyMarketBondHandler(w http.ResponseWriter, req *http.Request)
	SellMarketBondHandler(w http.ResponseWriter, req *http.Request)
}
