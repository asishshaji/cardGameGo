package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/asishshaji/go-voting-api/go/src/github.com/asishshaji/pitcherServer/models"
	"github.com/gorilla/mux"

	socketio "github.com/googollee/go-socket.io"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// App struct
type App struct {
	Router *mux.Router
	Client *mongo.Client
	DB     *mongo.Database
	l      *log.Logger
}

// Initialize function initializes the app
func (a *App) Initialize(dbname string) {

	connectionString := fmt.Sprintf("mongodb://localhost:27017/%s", dbname)

	var err error

	a.Client, err = mongo.NewClient(options.Client().ApplyURI(connectionString))

	if err != nil {
		a.l.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = a.Client.Connect(ctx)

	if err != nil {
		a.l.Fatal(err)
	}

	a.DB = a.Client.Database(dbname)

	defer a.Client.Disconnect(ctx)

	a.Router = mux.NewRouter()
	a.initializeRoutes()

}

func getCurrentTimeInNano() string {
	currentTime := time.Now().UnixNano()
	roomID := strconv.Itoa(int(currentTime))
	return roomID
}

// Room for the game
type Room struct {
	id      string
	sockets []socketio.Conn
}

func (a *App) initializeRoutes() {
	rooms := make(map[string]Room)

	socketServer, err := socketio.NewServer(nil)

	if err != nil {
		a.l.Fatal(err)
	}

	socketServer.OnConnect("/", func(s socketio.Conn) error {
		a.l.Println("New client connected " + s.ID())
		return nil
	})

	socketServer.OnEvent("/", "hostCreateNewGame", func(s socketio.Conn, msg string) {
		a.l.Println("Creating new room")
		room := Room{
			id: getCurrentTimeInNano(),
		}
		rooms[room.id] = room
		a.l.Println(room)

		// s.Emit("RoomID", room.id)

	})

	socketServer.OnEvent("/", "joinRoom", func(s socketio.Conn, msg string) {
		roomID := msg
		room := rooms[roomID]

		connectionsLen := len(room.sockets)
		if connectionsLen < 2 {
			room.sockets = append(room.sockets, s)
		}

	})

	socketServer.OnEvent("/", "getRooms", func(s socketio.Conn, msg string) {
		a.l.Println(s.Rooms())
		s.Emit("rooms", s.Rooms())
	})

	go socketServer.Serve()

	a.Router.Handle("/socket.io/", socketServer)

	// creates new card
	a.Router.HandleFunc("/card", a.createCard).Methods(http.MethodPost)

	// Create end point for generating unique link for a user
	// allow new users to join room via the link

}

func joinRoom(s socketio.Conn, room Room) {

}

func (a *App) createCard(rw http.ResponseWriter, r *http.Request) {

	rw.Header().Set("Content-Type", "application/json")

	var card models.Card

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&card)

	if err != nil {
		a.l.Println(err)
		http.Error(rw, "Error parsing card", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	err = card.CreateCard(a.DB)

	if err != nil {
		a.l.Fatalln(err)
	}

}

// Run starts the server
func (a *App) Run(port string) {

	server := &http.Server{
		Addr:    port,
		Handler: a.Router,
	}

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.l.Fatalf("listen: %s\n", err)
		}
	}()

	a.l.Printf("Server Started at PORT %v", server.Addr)

	// Until any cancellation signal is
	// received the code is blocked
	<-done
	a.l.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		a.l.Fatalf("Server Shutdown Failed:%+v", err)
	}
	a.l.Print("Server Shutdown Gracefully")
}
