package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"os"
)

type Server struct {
	Router *mux.Router
	DB     *sql.DB
}

func NewServer(db *sql.DB) *Server {
	server := Server{}

	server.DB = db
	server.Router = mux.NewRouter()
	return &server
}

func (s *Server) Start() error {
	return http.ListenAndServe(getListenAddr(), s.Router)
}

func getListenAddr() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	return net.JoinHostPort("", port)
}
