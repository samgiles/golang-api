package main

import (
    "errors"
	"encoding/json"
	"database/sql"
	"database/sql/driver"
	"github.com/lib/pq"
	"log"
)

type PostgresPaymentStore struct {
	db *sql.DB
}

func NewPostgresPaymentStore(db *sql.DB) *PostgresPaymentStore {
	return &PostgresPaymentStore{db}
}

func (s *PostgresPaymentStore) GetPayment(id string) (*Payment, error) {
	query := "SELECT (version, organisation_id, attributes) FROM payments WHERE id = $1 LIMIT 1;"

	p := Payment{Id: id}
	err := s.db.QueryRow(query, id).Scan(&p.Version, &p.OrganisationId, &p.Attributes)

	if err != nil {
        switch err {
        case sql.ErrNoRows:
            return nil, NewNotFoundError(err.Error())
        default:
		    log.Printf("Error getting payment: %s", err.Error())
            return nil, err
        }
	}

	return &p, nil
}

func (s *PostgresPaymentStore) GetAllPayments(page, limit int64) ([]Payment, error) {
	emptyList := make([]Payment, 0)
	return emptyList, nil
}

func (s *PostgresPaymentStore) CreatePayment(payment *Payment) (*Payment, error) {
    return s.create(payment)
}

func (s *PostgresPaymentStore) create(payment *Payment) (*Payment, error) {
    query := "SELECT id FROM InsertPaymentIdempotent($1,$2,$3,$4,$5);"

	err := s.db.QueryRow(query, payment.Id, payment.IdempotencyKey, payment.Version, payment.OrganisationId, payment.Attributes).Scan(&payment.Id)

	if err != nil {
		log.Printf("Error creating payment: %s", err.Error())

        pqerr, ok := err.(*pq.Error)

        if ok {
            return nil, convertPqError(pqerr)
        }

		return nil, err
	}

	return payment, nil
}

func (s *PostgresPaymentStore) UpdatePayment(payment *Payment) (*Payment, error) {
	return nil, NewNotFoundError("not found")
}

func (s *PostgresPaymentStore) DeletePayment(id string) error {
	return NewNotFoundError("not found")
}

// Implement database/driver.Valuer and sql.Scanner interfaces for attributes
// so we can easily scan them out of the db responses
func (attr PaymentAttributes) Value() (driver.Value, error) {
	return json.Marshal(attr)
}

func (attr *PaymentAttributes) Scan(src interface{}) error {
	source, typeValid := src.([]byte)

	if !typeValid {
		return errors.New("Type assertion failed: .([]byte) required")
	}

	return json.Unmarshal(source, attr)
}

func convertPqError(err *pq.Error) error {
    log.Printf("PQ ERROR CODE: %s", err.Code)
    switch err.Code {
    case "23505":
        return NewDocumentConflictError("")
    default:
        return err
    }
}
