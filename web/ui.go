package main

import (
	"encoding/json"
	"net/http"
)

func HandleUIRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	if r.Method != http.MethodGet {
		errorJson, _ := json.Marshal(Error{
			Message: "Method not allowed",
			Status:  http.StatusMethodNotAllowed,
		})
		http.Error(w, string(errorJson), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/ui" {
		http.ServeFile(w, r, "./ui/index.html")
		return
	}

	if r.URL.Path == "/js/app.js" {
		http.ServeFile(w, r, "./ui/js/app.js")
		return
	}

	if r.URL.Path == "/css/style.css" {
		http.ServeFile(w, r, "./ui/css/style.css")
		return
	}

	errorJson, _ := json.Marshal(Error{
		Message: "Not found",
		Status:  http.StatusNotFound,
	})

	http.Error(w, string(errorJson), http.StatusNotFound)
}
