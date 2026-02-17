CREATE TABLE dead_letter_notifications (
	id UUID PRIMARY KEY,
	notification_id UUID NOT NULL,
	payload JSONB NOT NULL,
	error TEXT NOT NULL,
	attempts INT NOT NULL,
	failed_at TIMESTAMP NOT NULL
);
