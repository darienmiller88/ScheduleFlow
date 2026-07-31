CREATE TABLE IF NOT EXISTS email_verification(
    id         SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    specialist_id INTEGER NOT NULL,
    expires_at    TIMESTAMPTZ NOT NULL,
    code_hash     TEXT NOT NULL,

    FOREIGN KEY (specialist_id) REFERENCES specialists(id) ON DELETE CASCADE
)