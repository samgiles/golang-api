package main

import (
	"database/sql"
	"database/sql/driver"
    "encoding/json"
	"errors"
	"github.com/lib/pq"
	"log"
    "time"

    "github.com/samgiles/health"
)

type PostgresPaymentStore struct {
    health.DefaultHealthCheck
	db *sql.DB
}

func NewPostgresPaymentStore(db *sql.DB) *PostgresPaymentStore {
    return &PostgresPaymentStore{db: db}
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
	// Use a CTE (common table expression) to capture the result of the update
	// so we can check whether we updated the id with the expected previous
	// version. This allows us to easily check the version update and whether
	// the id even exists in one statement. with the default iso level, READ
	// COMMITTED we won't overwrite data in the update thanks to the where
	// clause checking a committed version.
	query := `WITH updated_record AS ( UPDATE payments SET (version, organisation_id, attributes) = ($1, $2, $3) WHERE id = $4 AND version = $5 RETURNING *)
    (SELECT id, version, organisation_id, attributes FROM payments WHERE id = $4 UNION SELECT id, version, organisation_id, attributes FROM updated_record) LIMIT 2;`

	// TODO: Maybe versions should be opaque and non-sequential to avoid the
	// possibility of a client forcing an update by guessing the current db
	// version
	newVersion := payment.Version + 1
	rows, err := s.db.Query(query, newVersion, payment.OrganisationId, payment.Attributes, payment.Id, payment.Version)

	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		log.Printf("db_paymentstore: UPDATE ID NOT FOUND %s", payment.Id)
		return nil, NewNotFoundError("id not found")
	default:
		return nil, err
	}

	defer rows.Close()

	var latestPayment Payment

	for rows.Next() {
		if err := rows.Scan(&latestPayment.Id, &latestPayment.Version, &latestPayment.OrganisationId, &latestPayment.Attributes); err != nil {
			return nil, err
		}
	}

	// Return the DBs latest version, but also send a conflict error
	if latestPayment.Version != newVersion {
		log.Printf("db_paymentstore: update payment version mismatch expected version %d, actual %d", newVersion, latestPayment.Version)
		return &latestPayment, NewDocumentConflictError("conflict")
	}

	log.Printf("Updated to version %d", latestPayment.Version)

	// Return the updated latest version
	return &latestPayment, nil
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

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return NewNotFoundError("payment not found")
	}

	return nil
}

func (s *PostgresPaymentStore) RunHealthCheck() health.HealthCheckResult {
    // Using db.Ping is not enough, because they don't enter an error state
    // after being opened until someone tries to execute something on that
    // connection. With a healthcheck we'd prefer to be proactive
    if _, err := s.GetAllPayments(); err != nil {
        return health.NotReadyResult(err.Error())
    }

    return health.HealthyResult()
}

func (s *PostgresPaymentStore) HealthCheckName() string {
    return "postgresdb-health"
}

func (s *PostgresPaymentStore) HealthCheckFrequency() time.Duration {
    return 5 * time.Second
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
