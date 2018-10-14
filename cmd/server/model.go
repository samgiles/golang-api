package main

import (
	"encoding/json"
	"io"
)

type MonetaryAmount struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type ChargesInfo struct {
	BearerCode              string           `json:"bearer_code"`
	ReceiverChargesAmount   string           `json:"receiver_charges_amount"`
	ReceiverChargesCurrency string           `json:"receiver_charges_currency"`
	SenderCharges           []MonetaryAmount `json:"sender_charges"`
}

type PaymentParty struct {
	AccountName       string `json:"account_name,omitempty"`
	AccountNumber     string `json:"account_number,omitempty"`
	AccountNumberCode string `json:"account_number_code,omitempty"`
	AccountType       int32  `json:"account_type,omitempty"`
	Address           string `json:"address,omitempty"`
	BankId            string `json:"bank_id,omitempty"`
	BankIdCode        string `json:"bank_id_code,omitempty"`
	Name              string `json:"name,omitempty"`
}

type PaymentAttributes struct {
	Amount               string       `json:"amount"`
	Currency             string       `json:"currency"`
	EndToEndRef          string       `json:"end_to_end_reference"`
	NumericRef           string       `json:"numeric_reference"`
	PaymentId            string       `json:"payment_id"`
	PaymentPurpose       string       `json:"payment_purpose"`
	PaymentScheme        string       `json:"payment_scheme"`
	PaymentType          string       `json:"payment_type"`
	ProcessingDate       string       `json:"processing_date"`
	Reference            string       `json:"reference"`
	SchemePaymentType    string       `json:"scheme_payment_type"`
	SchemePaymentSubType string       `json:"scheme_payment_sub_type"`
	ChargesInfo          ChargesInfo  `json:"charges_information"`
	BeneficiaryParty     PaymentParty `json:"beneficiary_party"`
	DebtorParty          PaymentParty `json:"debtor_party"`
	SponsorParty         PaymentParty `json:"sponsor_party"`
}

type Payment struct {
	Id             string            `json:"id,omitempty" db:"id"`
	OrganisationId string            `json:"organisation_id" db:"organisation_id"`
	Version        int64             `json:"version" db:"version"`
	IdempotencyKey string            `json:"-" db:"idempotency_key"`
	Attributes     PaymentAttributes `json:"attributes" db:"attributes"`
}

type ErrorResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func NewErrorResponse(message string, status int) ErrorResponse {
	return ErrorResponse{Message: message, StatusCode: status}
}

type PaymentStore interface {
	// Returns the payment with specified id.  If the payment does not exist in
	// the store, `NotFoundError` is returned.
	// error may also contain unexpected db driver errors.
	GetPayment(id string) (*Payment, error)

	// Returns a list of all payments.
	// error may ontain unexpected db driver errors
	GetAllPayments() ([]Payment, error)

	// Creates a payment.  If payment.IdempotencyKey is set then that key will
	// be used as a surrogate key to prevent multiple inserts of the same
	// resource. If the resource already exists in the store, then a
	// `DocumentConflictError` is returned in error.
	//
	// error may also contain unexpected db driver errors.
	CreatePayment(payment *Payment) (*Payment, error)

	// Updates a resource at payment.Id. If payment.Version is not equal to the
	// current version held in the store then a `DocumentConflictError` is
	// returned, but also the latest payment resource in the database for that
	// Id is also returned in the Payment field.  These should be `nil`
	// checked by callers.
	// If the payment is not found, then a `NotFoundError` will be in the error
	// response
	// error may also contain unexpected db driver errors
	UpdatePayment(payment *Payment) (*Payment, error)

	// Deletes a payment from the store with specified id.
	// If the id is not found then error will be `NotFoundError`. If
	// successful, error will be `nil`.
	// error may also contain unexpected db driver errors.
	DeletePayment(id string) error
}

// Implements an "empty" payment store
type EmptyPaymentStore struct{}

func (s *EmptyPaymentStore) GetPayment(id string) (*Payment, error) {
	return nil, nil
}

func (s *EmptyPaymentStore) GetAllPayments() ([]Payment, error) {
	emptyList := make([]Payment, 0)
	return emptyList, nil
}

func (s *EmptyPaymentStore) CreatePayment(payment *Payment) (*Payment, error) {
	return nil, NewNotFoundError("not found")
}

func (s *EmptyPaymentStore) UpdatePayment(payment *Payment) (*Payment, error) {
	return nil, NewNotFoundError("not found")
}

func (s *EmptyPaymentStore) DeletePayment(id string) error {
	return NewNotFoundError("not found")
}

func UnmarshalPayment(data io.ReadCloser) (Payment, error) {
	var payment Payment
	decoder := json.NewDecoder(data)

	if err := decoder.Decode(&payment); err != nil {
		return payment, err
	}

	// TODO: Validate request

	return payment, nil
}
