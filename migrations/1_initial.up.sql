CREATE TABLE IF NOT EXISTS payments (
    id TEXT PRIMARY KEY,
    idempotency_key TEXT NOT NULL
        -- We don't want the user entering a HUGE key
        CHECK (char_length(idempotency_key) <= 100),
    version BIGINT NOT NULL,
    organisation_id TEXT NOT NULL,
        CHECK (char_length(idempotency_key) <= 100),
    attributes JSONB NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idempotency_keys_idx ON payments (idempotency_key);

CREATE OR REPLACE FUNCTION InsertPaymentIdempotent(
    _id payments.id%TYPE,
    _idempotency_key payments.idempotency_key%TYPE,
    _version payments.version%TYPE,
    _organisation_id payments.organisation_id%TYPE,
    _attributes payments.attributes%TYPE
)
RETURNS payments AS $$
    WITH new_row AS (
        INSERT INTO payments (id, idempotency_key, version, organisation_id, attributes)
        VALUES (_id, _idempotency_key, _version, _organisation_id, _attributes)
        ON CONFLICT (idempotency_key) DO NOTHING
        RETURNING *
    )
    SELECT * FROM new_row
    UNION
    SELECT * FROM payments WHERE idempotency_key = _idempotency_key;
$$ LANGUAGE sql;
