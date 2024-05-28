package brokering

type SendStatus int

const (
	SendFailed SendStatus = iota
	SendPanicked
	Sent
)

func (status SendStatus) String() string {
	return [...]string{"SendFailed", "SendPanicked", "Sent"}[status]
}
