package api

import (
	"net/http"
)

const DateFormat = "20060102"

func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", nextDayHandler)
}
