CREATE TABLE IF NOT EXISTS payments (
    id TEXT PRIMARY KEY,
    idempotency_key TEXT NOT NULL
        -- We don't want the user entering a HUGE key
        CHECK (char_length(idempotency_key) <= 100),
    version BIGINT NOT NULL,
    organisation_id TEXT NOT NULL,
        CHECK (char_length(idempotency_key) <= 100),
    payload JSONB NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idempotency_keys_idx ON payments (idempotency_key);
