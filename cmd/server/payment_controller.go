package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type PaymentController struct {
	store PaymentStore
}

func NewPaymentController(store PaymentStore) PaymentController {
	return PaymentController{store}
}

func (c *PaymentController) SetupRoutes(router *mux.Router) {
	router.HandleFunc("/payments", c.CreatePayment).Methods("POST")
	router.HandleFunc("/payments", c.ListPayments).Methods("GET")
	router.HandleFunc("/payments/{id}", c.GetPayment).Methods("GET")
	router.HandleFunc("/payments/{id}", c.UpdatePayment).Methods("PUT")
	router.HandleFunc("/payments/{id}", c.DeletePayment).Methods("DELETE")
}

func (c *PaymentController) CreatePayment(w http.ResponseWriter, r *http.Request) {

	payment, err := UnmarshalPayment(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var idemKey string
	if key := r.Header.Get("X-Idempotency-Key"); key != "" {
		idemKey = key
	} else {
		idemKey = uuid.New().String()
	}

	createdPayment, createErr := c.store.CreatePayment(payment, idemKey)

	if createErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(createErr.Error()))
		return
	}

	jsonResponse, _ := json.Marshal(createdPayment)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}

func (c *PaymentController) ListPayments(w http.ResponseWriter, r *http.Request) {
	payments, err := c.store.GetAllPayments(0, 100)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	jsonResponse, _ := json.Marshal(payments)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (c *PaymentController) GetPayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	payment, exists, _ := c.store.GetPayment(id)

	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonResponse, _ := json.Marshal(payment)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (c *PaymentController) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	payment, err := UnmarshalPayment(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	updatedPayment, err := c.store.UpdatePayment(id, payment.Version, payment)

	if err != nil {
		switch err.(type) {
		case *DocumentConflictError:
			w.WriteHeader(http.StatusConflict)
		case *NotFoundError:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Write([]byte(err.Error()))
		return
	}

	jsonResponse, _ := json.Marshal(updatedPayment)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (c *PaymentController) DeletePayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := c.store.DeletePayment(id)

	if err != nil {
		switch err.(type) {
		case *NotFoundError:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
