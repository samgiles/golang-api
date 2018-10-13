package main

import (
    "database/sql"
)

type PostgresPaymentStore struct {
    db *sql.DB
}

func (s *PostgresPaymentStore) GetPayment(id string) (*Payment, bool, error) {
	return nil, false, nil
}

func (s *PostgresPaymentStore) GetAllPayments(page, limit int64) ([]Payment, error) {
	emptyList := make([]Payment, 0)
	return emptyList, nil
}

func (s *PostgresPaymentStore) CreatePayment(payment Payment, idempotencyKey string) (*Payment, error) {
	return nil, NewNotFoundError("not found")
}

func (s *PostgresPaymentStore) UpdatePayment(id string, version int64, payment Payment) (*Payment, error) {
	return nil, NewNotFoundError("not found")
}

func (s *PostgresPaymentStore) DeletePayment(id string) error {
	return NewNotFoundError("not found")
}
