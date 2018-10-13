package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"os"

    "github.com/gorilla/handlers"
    "github.com/samgiles/health"
)

type Server struct {
	Router *mux.Router
	DB     *sql.DB

	paymentController PaymentController
    healthController health.HealthCheckController
}

func NewServer(db *sql.DB) *Server {
	server := Server{}

	server.DB = db
	server.Router = mux.NewRouter()

	server.paymentController = NewPaymentController(NewPostgresPaymentStore(db))
	server.paymentController.SetupRoutes(server.Router)
	server.healthController = health.NewHealthCheckController()

	return &server
}

func (s *Server) Start() error {
	return http.ListenAndServe(getListenAddr(), handlers.CombinedLoggingHandler(os.Stdout, s.Router))
}

func (s *Server) Stop() {
    s.healthController.Stop()
}

func getListenAddr() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	return net.JoinHostPort("", port)
}
