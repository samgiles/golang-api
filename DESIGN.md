# Payments API

## `POST /v1/payments`

Body must be `application/json` with the following required fields:

| Field         | Description           | Example  |
|:------------- |:-------------|:-----|
| `organisation_id` | `String` Organisation id of the creator of the Payment      | `organisation_id: '743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb'` |
| `attributes` | `Object` Payment attributes object      ||

(I haven't listed possible payment attributes in the interest of time :) )

Additionally an `X-Idempotency-Key` can be added to the header.  This allows
the client to retry in the case of unexpected errors safe in the knowledge
that the payment will only be created in the API once. The key should be a
random generated value, a uuid-v4 is suggested.

Will respond with HTTP status 201 (created) if successful, the response body
will contain a representation of the payment, including the assigned `id` and
`version` fields.

If the body is invalid, the API will respond with HTTP status 400 (bad
request).

## `GET /v1/payments`

Will list a set of payments.  Pagination is not implemented in this iteration,
it would be cool attempt to implement cursor based pagination as it would
probably scale better than a `LIMIT` based implementation.

## `GET /v1/payments/{id}`

If `{id}` not found, returns 404. Otherwise returns the entity as JSON.

## `PUT /v1/payments/{id}`

If `{id}` not found, returns 404.

This JSON payload must contain a `version` field. This `version` field should
remain the identical to the version of the resource you are updating.

For example:

If you get a resource from the API at `version: 6`, when making an update, you
should leave the `version` as `6`. This tells the API, I am updating version 6.
If when the update is sent to the database, version 6 was already updated, the
API will respond with HTTP status code 409 (conflict) and will respond with the
very latest version as a convenience - to avoid the additionally round trip to
get the updated version.

It is up to the client then to resolve concurrent update conflicts.

If the update was successful then will respond with http status code 200.

If the payload has an `id` field it is ignored. You can not update the id of a
resource.

## `DELETE /v1/payments/{id}`

If `{id}` not found, returns 404.
Else return http status code 204 (no content) as success indication.

# Implementation

Note: The requirements specified in this API design should be captured in
`main_test.go`.

The `PaymentController` accepts a `PaymentStore` interface that implements the
CRUD operations:

```GO
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
```

The payment controller then handles the HTTP requests, calls the appropriate
method from the store, and returns the appropriate error codes or success
response.

The semantics of this API are defined in `model.go` rather than duplicated here
in the design, and support the API design details listed above, including idempotent creates,
and versioned updates.

For simplicity we will use the `Payment` object everywhere as the DTO.

Postgres gives us ACID transactions and constraints that allow us to support
concurrent writes with versioning so that we don't experience lost updates as the API design
requires . We could use another data store that has this concept built in (couchdb for example),
but it's quite simple with postgres.

Additionally to support productionisation we will add healthcheck endpoints for
liveness and readiness (kubernetes' definitions).

## Out of scope things that would be nice to have in production

Opaque versions - ideally resource versions would be opaque rather than
sequential to make it impossible to guess the sequence and therefore try to
force updates.

Circuit breakers around the database calls.  Ideally we would circuit break
around a failing database to protect the database. The circuit state could tie
into the healthcheck.

Cursor based pagination for `GetAllPayments`.  Using `LIMIT <count> OFFSET
<offset>` doesn’t scale well for large datasets.  As offset increases the db
still needs to read `offset + count` rows before discarding the `offset` rows.
And if items are written frequently, and `GetAllPayments` always shows the most
recent additions, the page window becomes unreliable.  To implement cursor
based pagination we'd need a sort key in the DB. If this was production, it
would be a _must_ if data was going to grow to significant sizes.

Metrics, metrics metrics! We wouldn't know anything about our API without
observability.

Authentication, and then idempotency keys that are temporary and associated with
your auth. In this design they are global. If 10 years down the line you
generated an identical idempotency key, you'd lose your update unless you
regenerated. Unlikely with uuid idempotency keys, but technically possible.
