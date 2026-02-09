package modal

import "time"

type Status string

const (
	StatusPending    Status = "Pending"
	StatusProcessing Status = "Processing"
	StatusSend       Status = "Send"
	StatusFaild      Status = "Faild"
)

type Notification struct {
	Id        string
	Type      string
	Recipient string
	Payload   []byte
	Status    Status
	Attempts  int
	NextRetry *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
