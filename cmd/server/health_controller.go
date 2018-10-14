package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/samgiles/health"
)

type HealthController struct {
	hc health.HealthCheckController
}

func NewHealthController(hc health.HealthCheckController) HealthController {
	return HealthController{hc}
}

func (c *HealthController) SetupRoutes(router *mux.Router) {
	router.HandleFunc("/__readiness", c.HandleReadiness).Methods("GET")
	router.HandleFunc("/__liveness", c.HandleLiveness).Methods("GET")
}

func (c *HealthController) HandleReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ready := c.hc.Readiness()

	if ready.Ok {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	writeJsonResponse(w, ready)
}

func (c *HealthController) HandleLiveness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ready := c.hc.Liveness()

	if ready.Ok {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	writeJsonResponse(w, ready)
}
