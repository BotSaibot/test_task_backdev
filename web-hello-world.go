package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
	// })

	// http.ListenAndServe(":80", nil)
	servMux := http.NewServeMux()

	servMux.HandleFunc("/route1/", route1Handler)
	servMux.HandleFunc("/route2/", route2Handler)

	http.ListenAndServe(":80", servMux)
}

func route1Handler(wrt http.ResponseWriter, req *http.Request) {
	type reqMsg struct {
		msg, par0 string `json:"text"`
	}
	msg := fmt.Sprintf("Route 1:\nHello, you've requested: %s\n", req.URL.Path)
	js, err := json.Marshal(reqMsg{msg: msg, par0: "v any"})
	if err != nil {
		http.Error(wrt, err.Error(), http.StatusInternalServerError)
		return
	}
	wrt.Header().Set("Content-Type", "application/json")
	wrt.Write(js)
}

func route2Handler(wrt http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(wrt, "Route 2:\nHello, you've requested: %s\n", req.URL.Path)
}
