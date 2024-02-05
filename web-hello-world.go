package main

import (
	"encoding/json"
	"fmt"
	"log"
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

	log.Fatal(http.ListenAndServe(":8941", servMux))
}

type reqMsg struct {
	Msg  string `json:"msg"`
	Par0 string `json:"parameter"`
}

func route1Handler(wrt http.ResponseWriter, req *http.Request) {
	msg := fmt.Sprintf("Route 1:\nHello, you've requested: %s\n", req.URL.Path)
	js, err := json.MarshalIndent(reqMsg{Msg: msg, Par0: "v any"}, "", "\t")
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
