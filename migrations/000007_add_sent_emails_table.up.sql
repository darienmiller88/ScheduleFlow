CREATE TABLE sent_emails (
    id         SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    file_path     text NOT NULL,       -- snapshot of what was actually sent
    specialist_id INT NOT NULL,
    FOREIGN KEY (specialist_id) REFERENCES specialists(id)
);

CREATE INDEX idx_sent_emails_specialist_id ON sent_emails(specialist_id);