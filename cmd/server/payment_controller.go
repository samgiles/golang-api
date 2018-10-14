package main

import (
	"encoding/json"
	"log"
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
	w.Header().Set("Content-Type", "application/json")
	payment, err := UnmarshalPayment(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJsonResponse(w, NewErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	var idemKey string
	if key := r.Header.Get("X-Idempotency-Key"); key != "" {
		idemKey = key
	} else {
		idemKey = uuid.New().String()
	}

	payment.Id = uuid.New().String()
	payment.IdempotencyKey = idemKey

	createdPayment, createErr := c.store.CreatePayment(&payment)

	if createErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJsonResponse(w, NewErrorResponse(err.Error(), http.StatusInternalServerError))
		log.Printf("err: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	writeJsonResponse(w, createdPayment)
}

func (c *PaymentController) ListPayments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	payments, err := c.store.GetAllPayments()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJsonResponse(w, NewErrorResponse(err.Error(), http.StatusInternalServerError))
		log.Printf("err: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	writeJsonResponse(w, payments)
}

func (c *PaymentController) GetPayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	payment, err := c.store.GetPayment(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		writeJsonResponse(w, NewErrorResponse(err.Error(), http.StatusNotFound))
		log.Printf("err: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	writeJsonResponse(w, payment)
}

func (c *PaymentController) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	payment, err := UnmarshalPayment(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJsonResponse(w, NewErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	payment.Id = id
	updatedPayment, err := c.store.UpdatePayment(&payment)

	if err != nil {
		var responseCode int
		switch err.(type) {
		case *DocumentConflictError:
			responseCode = http.StatusConflict
		case *NotFoundError:
			responseCode = http.StatusNotFound
		default:
			responseCode = http.StatusInternalServerError
		}

		w.WriteHeader(responseCode)
		writeJsonResponse(w, NewErrorResponse(err.Error(), responseCode))
		log.Printf("err: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	writeJsonResponse(w, updatedPayment)
}

func (c *PaymentController) DeletePayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	err := c.store.DeletePayment(id)

	if err != nil {
		var responseCode int
		switch err.(type) {
		case *NotFoundError:
			responseCode = http.StatusNotFound
		default:
			log.Printf("err: %s", err.Error())
			responseCode = http.StatusInternalServerError
		}

		w.WriteHeader(responseCode)
		writeJsonResponse(w, NewErrorResponse(err.Error(), responseCode))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJsonResponse(w http.ResponseWriter, obj interface{}) {
	jsonResponse, _ := json.Marshal(obj)
	w.Write(jsonResponse)
}
