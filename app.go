package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/asishshaji/go-voting-api/go/src/github.com/asishshaji/pitcherServer/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// App
type App struct {
	Router *mux.Router
	Client *mongo.Client
	DB     *mongo.Database
}

// Initialized app
func (a *App) Initialize(dbname string) {
	connectionString := fmt.Sprintf("mongodb://localhost:27017/%s", dbname)

	var err error

	a.Client, err = mongo.NewClient(options.Client().ApplyURI(connectionString))

	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = a.Client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	a.DB = a.Client.Database(dbname)

	defer a.Client.Disconnect(ctx)

	a.Router = mux.NewRouter()
	a.initializeRoutes()

}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/card", a.createCard).Methods(http.MethodPost)
}

func (a *App) createCard(rw http.ResponseWriter, r *http.Request) {

	var card models.Card

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&card); err != nil {
		http.Error(rw, "Error parsing card", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	_ = card.CreateCard(a.DB)

}

func (a *App) Run(addr string) {
	log.Fatalln(http.ListenAndServe(addr, a.Router))
}
