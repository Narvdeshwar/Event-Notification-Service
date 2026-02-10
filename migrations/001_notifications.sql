CREATE TABLE notifications{
    id UUID PRIMARY KEY,
    type TEXT NOT NULL,
    recipient TEXT NOT NULL,
    payload JSONB NOT NULL,
    status TEXT NOT NULL,
    attempts INT NOT NULL DEFAULT 0,
    next_retry_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
}

CREATE INDEX idx_notification_pending ON notifications(
    status,
    next_retry_at
)
