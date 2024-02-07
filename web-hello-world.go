package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"autorization_test/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	customHandler struct {
		database    *mongo.Client
		handlerFunc func(http.ResponseWriter, *http.Request, *mongo.Client)
	}
	User struct {
		Name string
		GUID string
	}
	reqMsg struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		DebugGUID    string `json:"debug-guid"`
	}
	session struct {
		User         User
		RefreshToken string
	}
)

func init() {
	err := config.Export("./config.json")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	database, err := mongo.NewClient(options.Client().ApplyURI(config.Get().DatabaseURL))
	if err != nil {
		log.Fatal(err)
	}

	err = database.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Connected to MongoDB!")

	servMux := http.NewServeMux()

	servMux.Handle("/route1/", customHandler{database, route1Handler})
	servMux.HandleFunc("/route2/", route2Handler)
	servMux.HandleFunc("/debug/", debugHandler)
	servMux.Handle("/allUsers/", customHandler{database, allUsersHandler})
	servMux.Handle("/addUser/", customHandler{database, addUserHandler})

	addr := fmt.Sprintf("%s:%s", config.Get().IP, config.Get().Port)
	log.Printf("Http server (%s) run...", addr)
	log.Fatal(http.ListenAndServe(addr, servMux))
}

func (handler customHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler.handlerFunc(w, req, handler.database)
}

func route1Handler(w http.ResponseWriter, req *http.Request, db *mongo.Client) {
	// Первый маршрут выдает пару Access, Refresh токенов для пользователя с идентификатором (GUID) указанным в параметре запроса

	if req.Method == http.MethodPost || req.Method == http.MethodGet {
		formGUID := req.FormValue("guid")
		secretKey := []byte(config.Get().Secret)
		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
			"guid": formGUID,
			"eat":  time.Now().Add(time.Hour).Unix(),
			"iat":  time.Now().Unix(),
		})
		SignedAccessToken, err := accessToken.SignedString(secretKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		newUUID := []byte(uuid.New().String())
		refreshToken := base64.StdEncoding.EncodeToString(newUUID)

		js, err := json.MarshalIndent(reqMsg{
			AccessToken:  SignedAccessToken,
			RefreshToken: refreshToken,
			DebugGUID:    formGUID,
		}, "", "\t")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		collection := db.Database("test_task_backend").Collection("sessions")
		addSession := session{
			User: User{
				GUID: formGUID,
			},
			RefreshToken: refreshToken,
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

func allUsersHandler(w http.ResponseWriter, req *http.Request, db *mongo.Client) {
	var results []interface{}
	fmt.Fprintf(w, "Route 'All users':\nHello, you've requested: %s.\n", req.URL.Path)
	collection := db.Database("test_task_backend").Collection("users")
	cur, err := collection.Find(context.TODO(), bson.M{}, options.Find())
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var el struct {
			User User
		}

		err := cur.Decode(&el)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, el.User)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	fmt.Fprintf(w, "Results:\n")
	for i, v := range results {
		fmt.Fprintf(w, "\t%v - %+v\n", i, v)
	}
}

func addUserHandler(w http.ResponseWriter, req *http.Request, db *mongo.Client) {
	fmt.Fprintf(w, "Route 'Add user':\nHello, you've requested: %s.\n", req.URL.Path)
	collection := db.Database("test_task_backend").Collection("users")
	// tmpInsert := bson.M{"user": bson.M{"name": "Bob", "guid": "89450g8"}}
	tmpInsert := bson.M{"user": User{"Test-user", "-8u3456yji"}}
	insertResult, err := collection.InsertOne(context.TODO(), tmpInsert)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, "Inserted a single document: ", insertResult.InsertedID)
}
