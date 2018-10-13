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

	paymentController PaymentController
}

func NewServer(db *sql.DB) *Server {
	server := Server{}

	server.DB = db
	server.Router = mux.NewRouter()

	server.paymentController = NewPaymentController(&EmptyPaymentStore{})
	server.paymentController.SetupRoutes(server.Router)

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
