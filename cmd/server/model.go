package main

import "io"
import "encoding/json"

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
	Id             string            `json:"id,omitempty"`
	OrganisationId string            `json:"organisation_id"`
	Version        int64             `json:"version"`
	Attributes     PaymentAttributes `json:"attributes"`
}

type PaymentStore interface {
	GetPayment(id string) (*Payment, bool, error)
	GetAllPayments(page, limit int64) ([]Payment, error)
	CreatePayment(payment Payment, idempotencyKey string) (*Payment, error)
	UpdatePayment(id string, version int64, payment Payment) (*Payment, error)
	DeletePayment(id string) error
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
