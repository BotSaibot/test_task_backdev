package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type SingIn struct {
	GUID string `json:"guid"`
}

func main() {
	servMux := http.NewServeMux()

	servMux.HandleFunc("/route1/", route1Handler)
	servMux.HandleFunc("/route2/", route2Handler)
	servMux.HandleFunc("/debug/", debugHandler)

	log.Fatal(http.ListenAndServe(":8941", servMux))
}

func route1Handler(w http.ResponseWriter, req *http.Request) {
	// Первый маршрут выдает пару Access, Refresh токенов для пользователя сидентификатором (GUID) указанным в параметре запроса
	type reqMsg struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		DebugGUID    string `json:"debug-guid"`
	}

	if req.Method == http.MethodPost {
		formVar := req.FormValue("guid")
		js, err := json.MarshalIndent(reqMsg{
			AccessToken:  "someAccessToken",
			RefreshToken: "RefreshToken",
			DebugGUID:    formVar,
		}, "", "\t")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}
}

func route2Handler(w http.ResponseWriter, req *http.Request) {
	// Второй маршрут выполняет Refresh операцию на пару Access, Refreshтокенов
	fmt.Fprintf(w, "Route 2:\nHello, you've requested: %s\n", req.URL.Path)
}

func debugHandler(w http.ResponseWriter, req *http.Request) {
	debPage := `<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title></title>
</head>
<body>
	<h1>Debug Route</h1>
	<p>Вот два маршрута: <a href="/route1/">1 маршрут</a>; <a href="/route2/">2 маршрут</a></p>
	<form action="/route1/" method="post">
		<p>GUID: <input type="text" name="guid"></p>
		<p><input type="submit"></p>
	</form>
</body>
</html>`

	fmt.Fprintf(w, debPage)
}
