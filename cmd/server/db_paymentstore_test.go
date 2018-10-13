package main

import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func TestGetPaymentNonExistent(t *testing.T) {
	paymentStore, mock, db := createPaymentStore(t)
	defer db.Close()

	mock.ExpectQuery("SELECT (.+) FROM payments WHERE id = (.+) LIMIT 1;").WillReturnError(sql.ErrNoRows)

	_, err := paymentStore.GetPayment("non-existent-id")

	assert.Error(t, NewNotFoundError(""), err)
	assertExpectations(t, mock)
}

func TestGetPaymentExists(t *testing.T) {
	paymentStore, mock, db := createPaymentStore(t)
	defer db.Close()

	rows := [3]string{"version", "organisation_id", "attributes"}
	queryResult := sqlmock.NewRows(rows[:])

	attrsVal, _ := createPaymentAttributes("100").Value()
	// Call the Valuer (Value method) created for payment attrs
	queryResult = queryResult.AddRow(0, "organisation_id", attrsVal)

	mock.ExpectQuery("SELECT (.+) FROM payments WHERE id = (.+) LIMIT 1;").WillReturnRows(queryResult)

	payment, err := paymentStore.GetPayment("anid")

	if err != nil {
		t.Fatalf("Unexpected error from payment store: %s", err)
	}

	assert.Equal(t, payment.Id, "anid")
	assert.Equal(t, payment.Version, int64(0))
	assert.Equal(t, payment.OrganisationId, "organisation_id")
	assert.Equal(t, payment.Attributes.Amount, "100")
	assertExpectations(t, mock)
}

func TestGetAllPaymentsNoRows(t *testing.T) {
	paymentStore, mock, db := createPaymentStore(t)
	defer db.Close()

	mock.ExpectQuery("SELECT (.+) FROM payments LIMIT (.+);").WillReturnError(sql.ErrNoRows)

	payments, err := paymentStore.GetAllPayments()

	if err != nil {
		t.Fatalf("unexpected error fetching rows: %s", err)
	}

	assert.Equal(t, 0, len(payments))
	assertExpectations(t, mock)
}

func TestGetAllPaymentsWithRows(t *testing.T) {
	paymentStore, mock, db := createPaymentStore(t)
	defer db.Close()

	rows := [4]string{"id", "version", "organisation_id", "attributes"}
	queryResult := sqlmock.NewRows(rows[:])

	attrsValA, _ := createPaymentAttributes("200").Value()
	queryResult = queryResult.AddRow("anid1", 0, "organisation_id", attrsValA)

	attrsValB, _ := createPaymentAttributes("300").Value()
	queryResult = queryResult.AddRow("anid2", 0, "organisation_id", attrsValB)

	mock.ExpectQuery("SELECT (.+) FROM payments LIMIT (.+);").WillReturnRows(queryResult)

	payments, err := paymentStore.GetAllPayments()

	if err != nil {
		t.Fatalf("unexpected error fetching rows: %s", err)
	}

	assert.Equal(t, 2, len(payments))
	assertExpectations(t, mock)
}

func TestUpdatePaymentNonExistent(t *testing.T) {
	paymentStore, mock, db := createPaymentStore(t)
	defer db.Close()

	mock.ExpectQuery("WITH updated_record AS").WillReturnError(sql.ErrNoRows)

	_, err := paymentStore.UpdatePayment(&Payment{Id: "id", Version: 0})
	assert.Error(t, NewNotFoundError(""), err)
	assertExpectations(t, mock)
}

func TestUpdatePaymentVersionConflictResponse(t *testing.T) {
	paymentStore, mock, db := createPaymentStore(t)
	defer db.Close()

	rows := [4]string{"id", "version", "organisation_id", "attributes"}
	queryResult := sqlmock.NewRows(rows[:])

	// We mock this as the current value in the database, it's at version 1
	attrsValOld, _ := createPaymentAttributes("200").Value()
	queryResult = queryResult.AddRow("anid", 1, "organisation_id", attrsValOld)

	mock.ExpectQuery("WITH updated_record AS").WillReturnRows(queryResult)

	// We update with something we thought was at version 0
	payment, err := paymentStore.UpdatePayment(&Payment{Id: "id", Version: 0})

	assert.Error(t, NewDocumentConflictError(""), err)
	assert.Equal(t, int64(1), payment.Version)
	assertExpectations(t, mock)
}

func TestUpdatePaymentNoConflict(t *testing.T) {
	paymentStore, mock, db := createPaymentStore(t)
	defer db.Close()

	rows := [4]string{"id", "version", "organisation_id", "attributes"}
	queryResult := sqlmock.NewRows(rows[:])

	// We mock this as the current value in the database, it's at version 1
	attrsValOld, _ := createPaymentAttributes("200").Value()
	queryResult = queryResult.AddRow("anid", 1, "organisation_id", attrsValOld)

	// We mock this as the current value in the database, it's at version 1
	attrsValNew, _ := createPaymentAttributes("300").Value()
	queryResult = queryResult.AddRow("anid", 2, "organisation_id", attrsValNew)

	mock.ExpectQuery("WITH updated_record AS").WillReturnRows(queryResult)

	// We update with something we though was at version 1
	payment, err := paymentStore.UpdatePayment(&Payment{Id: "id", Version: 1})

	if err != nil {
		t.Errorf("Unexpected error updating payment: %s", err)
	}

	assert.Equal(t, int64(2), payment.Version)
	assert.Equal(t, "300", payment.Attributes.Amount)
	assertExpectations(t, mock)
}

func TestCreatePayment(t *testing.T) {
	paymentStore, mock, db := createPaymentStore(t)
	defer db.Close()

	rows := [1]string{"id"}
	queryResult := sqlmock.NewRows(rows[:])
	queryResult = queryResult.AddRow("anid")

	mock.ExpectQuery("SELECT id FROM InsertPaymentIdempotent").WillReturnRows(queryResult)

	createdPayment, err := paymentStore.CreatePayment(&Payment{Id: "anid"})

	if err != nil {
		t.Errorf("Unexpected error creating payment: %s", err)
	}

	assert.Equal(t, "anid", createdPayment.Id)
	assertExpectations(t, mock)
}

func TestCreatePaymentIdempotencyKeyConstraintViolation(t *testing.T) {
	paymentStore, mock, db := createPaymentStore(t)
	defer db.Close()

	// pg error constant
	constraintViolationError := pq.ErrorCode("23505")
	mock.ExpectQuery("SELECT id FROM InsertPaymentIdempotent").WillReturnError(pq.Error{Code: constraintViolationError})

	_, err := paymentStore.CreatePayment(&Payment{Id: "anid", IdempotencyKey: "idemkey"})

	assert.Error(t, NewDocumentConflictError(""), err)
	assertExpectations(t, mock)
}

func TestDeletePaymentNonExistent(t *testing.T) {
	paymentStore, mock, db := createPaymentStore(t)
	defer db.Close()

	mock.ExpectExec("DELETE FROM payments WHERE id").WillReturnError(sql.ErrNoRows)

	err := paymentStore.DeletePayment("id")

	assert.Error(t, NewNotFoundError(""), err)
	assertExpectations(t, mock)
}

func TestDeletePaymentExists(t *testing.T) {
	paymentStore, mock, db := createPaymentStore(t)
	defer db.Close()

	mock.ExpectExec("DELETE FROM payments WHERE id").WillReturnResult(sqlmock.NewResult(-1, 1))

	err := paymentStore.DeletePayment("id")

	assert.Equal(t, nil, err)
	assertExpectations(t, mock)
}

func createPaymentStore(t *testing.T) (PaymentStore, sqlmock.Sqlmock, *sql.DB) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("Could not create mock: %s", err)
	}

	return NewPostgresPaymentStore(db), mock, db
}

func createPaymentAttributes(amount string) PaymentAttributes {
	attrs := PaymentAttributes{}
	attrs.Amount = amount
	attrs.Currency = "GBP"
	attrs.EndToEndRef = "e2eref"
	attrs.NumericRef = "12313123"
	attrs.PaymentId = "123"
	return attrs
}

func assertExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectation error: %s", err)
	}
}
