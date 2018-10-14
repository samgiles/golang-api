package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"log"
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

    paymentStore := NewPostgresPaymentStore(db)

	server.paymentController = NewPaymentController(paymentStore)
	server.paymentController.SetupRoutes(server.Router)

    healthchecks := health.NewHealthCheckController()
    healthchecks.AddHealthCheck(paymentStore)
    server.healthController = healthchecks

    healthController := NewHealthController(healthchecks)
    healthController.SetupRoutes(server.Router)

	return &server
}

func (s *Server) Start() error {
	addr := getListenAddr()
	log.Printf("server: starting http listener on %s", addr)
	return http.ListenAndServe(addr, handlers.CombinedLoggingHandler(os.Stdout, s.Router))
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
