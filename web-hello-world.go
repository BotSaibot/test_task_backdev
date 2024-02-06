package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	singIn struct {
		GUID string `json:"guid"`
	}
	customHandler struct {
		database *mongo.Client
	}
	User struct {
		Name string
		GUID string
	}
)

func main() {
	database, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/"))
	if err != nil {
		log.Fatal(err)
	}

	err = database.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	servMux := http.NewServeMux()

	servMux.HandleFunc("/route1/", route1Handler)
	servMux.HandleFunc("/route2/", route2Handler)
	servMux.HandleFunc("/debug/", debugHandler)
	servMux.Handle("/route2v2/", customHandler{database})

	log.Fatal(http.ListenAndServe(":8941", servMux))
}

func route1Handler(w http.ResponseWriter, req *http.Request) {
	// Первый маршрут выдает пару Access, Refresh токенов для пользователя с идентификатором (GUID) указанным в параметре запроса
	type reqMsg struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		DebugGUID    string `json:"debug-guid"`
	}

	if req.Method == http.MethodPost || req.Method == http.MethodGet {
		formGUID := req.FormValue("guid")
		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
			"guid": formGUID,
			"eat":  time.Now().Add(time.Hour).Unix(),
			"iat":  time.Now().Unix(),
		})
		SignedAccessToken, err := accessToken.SignedString([]byte("secret"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		js, err := json.MarshalIndent(reqMsg{
			// AccessToken: "someAccessToken",
			AccessToken:  SignedAccessToken,
			RefreshToken: "RefreshToken",
			DebugGUID:    formGUID,
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
	// Второй маршрут выполняет Refresh операцию на пару Access, Refresh токенов
	fmt.Fprintf(w, "Route 2:\nHello, you've requested: %s\n", req.URL.Path)
	// data
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

func (handler customHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var results []bson.M
	fmt.Fprintf(w, "Route 2v2:\nHello, you've requested: %s\n", req.URL.Path)
	// handler.database
	collection := handler.database.Database("test_task_backend").Collection("users")
	// tmpInsert := bson.M{"user": bson.M{"name": "Bob", "guid": "89450g8"}}
	// tmpInsert := bson.M{"user": User{"Alice", "(*j5jsd@"}}
	// insertResult, err := collection.InsertOne(context.TODO(), tmpInsert)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Fprint(w, "Inserted a single document: ", insertResult.InsertedID)
	cur, err := collection.Find(context.TODO(), bson.M{}, options.Find())
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var el bson.M
		err := cur.Decode(&el)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, el)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	fmt.Fprintf(w, "Results: %+v\n", results)
}
