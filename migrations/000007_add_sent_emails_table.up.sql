CREATE TABLE sent_emails (
    id         SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    specialist_id INT NOT NULL REFERENCES specialists(specialist_id),
    file_path     text NOT NULL,       -- snapshot of what was actually sent
);

CREATE INDEX idx_sent_emails_specialist_id ON sent_emails(specialist_id);