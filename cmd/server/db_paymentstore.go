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
	query := "SELECT version, organisation_id, attributes FROM payments WHERE id = $1 LIMIT 1;"

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

func (s *PostgresPaymentStore) GetAllPayments() ([]Payment, error) {
    query := "SELECT id, version, organisation_id, attributes FROM payments LIMIT $1;"

    // TODO: offset/limit based pagination - or a more scalable forward cursor based pagination?
    rows, err := s.db.Query(query, 100)

    switch err {
    case nil:
        break
    case sql.ErrNoRows:
        return make([]Payment, 0), nil
    default:
        return nil, err
    }

    defer rows.Close()

    payments := make([]Payment, 0)

    for rows.Next() {
        payment := Payment{}
        if err := rows.Scan(&payment.Id, &payment.Version, &payment.OrganisationId, &payment.Attributes); err != nil {
            return nil, err
        }

        payments = append(payments, payment)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

	return payments, nil
}

func (s *PostgresPaymentStore) CreatePayment(payment *Payment) (*Payment, error) {
    query := "SELECT id FROM InsertPaymentIdempotent($1,$2,$3,$4,$5);"

	err := s.db.QueryRow(query, payment.Id, payment.IdempotencyKey, payment.Version, payment.OrganisationId, payment.Attributes).Scan(&payment.Id)

	if err != nil {
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
    query := "DELETE FROM payments WHERE id = $1;"
	result, err := s.db.Exec(query, id)

    if err != nil {
        switch err {
        case sql.ErrNoRows:
            return NewNotFoundError(err.Error())
        default:
            return err
        }
	}

    rows, err := result.RowsAffected();
    if err != nil {
        return err
    }

    if rows == 0 {
        return NewNotFoundError("payment not found")
    }

	return nil
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
    switch err.Code {
    case "23505":
        return NewDocumentConflictError("")
    default:
        return err
    }
}
