package model

import "time"

type Status string

const (
	StatusPending    Status = "Pending"
	StatusProcessing Status = "Processing"
	StatusSent       Status = "Sent"
	StatusFailed     Status = "Failed"
)

type Notification struct {
	ID        string
	Type      string
	Recipient string
	Payload   []byte
	Status    Status
	Attempts  int
	NextRetry *time.Time
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
