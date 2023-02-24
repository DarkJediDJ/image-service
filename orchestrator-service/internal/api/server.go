package server

import (
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"

	http_handler "github.com/DarkJediDJ/image-service/orchestrator-service/internal/api/http"
)

type App struct {
	Router *mux.Router
}

func NewApp() *App {
	return &App{}
}

func (a *App) InitRouter(conn *kafka.Conn, session *session.Session) {
	myRouter := mux.NewRouter().StrictSlash(false)
	myRouter.HandleFunc("/", http_handler.NewHandler(conn, session).Process).Methods("POST")
	a.Router = myRouter
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
