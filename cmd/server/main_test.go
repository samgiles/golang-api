package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testOrgId = "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb"

var server *Server

func setUpApplication() error {
	testDbName := os.Getenv("TEST_DB_NAME")
	db, err := CreateDatabaseConnection(
		os.Getenv("TEST_DB_USERNAME"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_HOST"),
		testDbName)

	if err != nil {
		return err
	}

	server = NewServer(db)

	log.Println("Starting DB connection wait")
	err = WaitForDbConnectivity(db, 10*time.Second)

	if err != nil {
		return err
	}

	log.Println("Connected to DB successfully")

	log.Println("Migrating database up..")
	err = MigrateDatabaseUp(testDbName, db)

	if err != nil {
		return err
	}

	log.Println("Migrated database up")
	return nil
}

func clearDb() {
    server.DB.Exec("TRUNCATE TABLE payments;")
}

func TestApplicationIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	if err := setUpApplication(); err != nil {
		t.Errorf("Failed to start application under test: %s", err)
		return
	}

	t.Run("POST /payments", func(t *testing.T) {
		defer clearDb()

		t.Run("An invalid resource to should fail with bad request status", func(t *testing.T) {
			testPayment := []byte(" \"invalid\": 0 }")

			req, _ := http.NewRequest("POST", "/payments", bytes.NewBuffer(testPayment))

			response := executeRequest(req)
			checkResponseCode(t, http.StatusBadRequest, response.Code)
		})

		t.Run("A valid resource should create a new payment and respond with the new resource", func(t *testing.T) {
			// Additionally:
			//  - should allocate a resource ID
			//  - should set the resource's version to 0

			payment := createNewPayment(t)

			assert.NotEqual(t, "", payment.Id, "Expected payment.Id but was empty")
			assert.Equal(t, testOrgId, payment.OrganisationId)
			assert.Equal(t, int64(0), payment.Version)
		})

		t.Run("Creating a resource with the same idempotency key should be idempotent", func(t *testing.T) {
			// Should allocate only one resource ID and always return the
			// same resource for the given idempotency key

			testPayment, e := ioutil.ReadFile("../../test/test_payment.json")

			if e != nil {
				t.Error(e)
				return
			}

			reqFirst, _ := http.NewRequest("POST", "/payments", bytes.NewBuffer(testPayment))
			reqFirst.Header.Add("X-Idempotency-Key", "abcd")

			reqSecond, _ := http.NewRequest("POST", "/payments", bytes.NewBuffer(testPayment))
			reqSecond.Header.Add("X-Idempotency-Key", "abcd")

			responseFirst := executeRequest(reqFirst)
			responseSecond := executeRequest(reqSecond)

			checkResponseCode(t, http.StatusCreated, responseFirst.Code)
			checkResponseCode(t, http.StatusCreated, responseSecond.Code)

			paymentFirst := readPaymentResponse(responseFirst.Body.Bytes())
			paymentSecond := readPaymentResponse(responseSecond.Body.Bytes())

			assert.NotEqual(t, "", paymentFirst.Id)
			assert.NotEqual(t, "", paymentSecond.Id)
			assert.Equal(t, testOrgId, paymentFirst.OrganisationId)
			assert.Equal(t, testOrgId, paymentSecond.OrganisationId)
			assert.Equal(t, paymentFirst.OrganisationId, paymentSecond.OrganisationId)
		})
	})

	t.Run("GET /payments", func(t *testing.T) {
		t.Run("Should return an empty array if no resources exist", func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/payments", nil)
			response := executeRequest(req)
			checkResponseCode(t, http.StatusOK, response.Code)

			assert.Equal(t, "[]", response.Body.String())
		})

		t.Run("Should list payment resources after their creation", func(t *testing.T) {
			defer clearDb()
			createNewPayment(t)

			req, _ := http.NewRequest("GET", "/payments", nil)

			responseGet := executeRequest(req)
			checkResponseCode(t, http.StatusOK, responseGet.Code)

			paymentList := readPaymentListResponse(responseGet.Body.Bytes())

			if assert.NotZero(t, len(paymentList)) {
				assert.Equal(t, paymentList[0].OrganisationId, testOrgId)
			}
		})
	})

	t.Run("GET /payments/{id}", func(t *testing.T) {
		t.Run("A missing payment resource should return not found status code", func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/payments/a-non-existent-id", nil)
			response := executeRequest(req)
			checkResponseCode(t, http.StatusNotFound, response.Code)
		})

		t.Run("Should return a resource after it was created", func(t *testing.T) {
			defer clearDb()
			paymentCreated := createNewPayment(t)

			if !assert.NotEqual(t, "", paymentCreated.Id) {
				return
			}

			paymentResource := getPayment(t, paymentCreated.Id)

			assert.NotEqual(t, "", paymentResource.Id)
			assert.Equal(t, paymentCreated.Id, paymentResource.Id)
		})
	})

	t.Run("DELETE /payments/{id}", func(t *testing.T) {
		t.Run("Deleting a missing payment should return not found status code", func(t *testing.T) {
			del, _ := http.NewRequest("DELETE", "/payments/non-existent-id", nil)
			responseDel := executeRequest(del)
			checkResponseCode(t, http.StatusNotFound, responseDel.Code)
		})

		t.Run("After deleting a payment, attempting to get resource should return not found status code", func(t *testing.T) {
			defer clearDb()
			paymentCreated := createNewPayment(t)

			if paymentCreated.Id == "" {
				t.Errorf("Created Id was empty")
				return
			}

			deletePayment(t, paymentCreated.Id)

			resourcePath := path.Join("/payments", paymentCreated.Id)
			del, _ := http.NewRequest("GET", resourcePath, nil)
			responseDel := executeRequest(del)
			checkResponseCode(t, http.StatusNotFound, responseDel.Code)
		})
	})

	t.Run("PUT /payments/{id}", func(t *testing.T) {
		t.Run("Updating a value on a resource should be persisted for subsequent gets of the resource", func(t *testing.T) {
			// Additionally
			//  - the newly updated resource should be returned after the
			//    update to the client
			defer clearDb()

			paymentCreated := createNewPayment(t)

			if paymentCreated.Id == "" {
				t.Errorf("Created Id was empty")
				return
			}

			// Update amount
			paymentCreated.Attributes.Amount = "2000"

			updatedPayment := updatePayment(t, paymentCreated)

			assert.NotEqual(t, "", updatedPayment.Id)
			assert.Equal(t, paymentCreated.Id, updatedPayment.Id)
			assert.Equal(t, "2000", updatedPayment.Attributes.Amount)

			readPayment := getPayment(t, paymentCreated.Id)
			assert.Equal(t, "2000", readPayment.Attributes.Amount)
		})

		t.Run("A conflict status code should be returned if the client version is behind the server's version", func(t *testing.T) {
			defer clearDb()
			paymentCreated := createNewPayment(t)

			if paymentCreated.Id == "" {
				t.Errorf("Created Id was empty")
				return
			}

			// Manually update the payment version locally to mock differing client
			// version
			paymentCreated.Version = 2
			paymentCreated.Attributes.Amount = "3000"

			resourcePath := path.Join("/payments", paymentCreated.Id)
			json, _ := json.Marshal(paymentCreated)
			putReq, _ := http.NewRequest("PUT", resourcePath, bytes.NewBuffer(json))
			putResponse := executeRequest(putReq)
			checkResponseCode(t, http.StatusConflict, putResponse.Code)
		})
	})
}

func createNewPayment(t *testing.T) Payment {
	testPayment, e := ioutil.ReadFile("../../test/test_payment.json")

	if e != nil {
		t.Error(e)
	}

	create, _ := http.NewRequest("POST", "/payments", bytes.NewBuffer(testPayment))
	responseCreate := executeRequest(create)
	checkResponseCode(t, http.StatusCreated, responseCreate.Code)

	return readPaymentResponse(responseCreate.Body.Bytes())
}

func getPayment(t *testing.T, id string) Payment {
	resourcePath := path.Join("/payments", id)
	get, _ := http.NewRequest("GET", resourcePath, nil)
	responseGet := executeRequest(get)
	checkResponseCode(t, http.StatusOK, responseGet.Code)

	return readPaymentResponse(responseGet.Body.Bytes())
}

func updatePayment(t *testing.T, payment Payment) Payment {
	resourcePath := path.Join("/payments", payment.Id)
	json, _ := json.Marshal(payment)
	putReq, _ := http.NewRequest("PUT", resourcePath, bytes.NewBuffer(json))
	putResponse := executeRequest(putReq)
	checkResponseCode(t, http.StatusOK, putResponse.Code)

	return readPaymentResponse(putResponse.Body.Bytes())
}

func deletePayment(t *testing.T, id string) {
	resourcePath := path.Join("/payments", id)
	del, _ := http.NewRequest("DELETE", resourcePath, nil)
	responseDel := executeRequest(del)
	checkResponseCode(t, http.StatusNoContent, responseDel.Code)
}

func readPaymentResponse(body []byte) Payment {
	var payment Payment
	json.Unmarshal(body, &payment)
	return payment
}

func readPaymentListResponse(body []byte) []Payment {
	var payments []Payment
	json.Unmarshal(body, &payments)
	return payments
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	server.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Was %d\n", expected, actual)
	}
}
