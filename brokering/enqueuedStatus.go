package brokering

type EnqueuedStatus int

const (
	EnqueuingFailed EnqueuedStatus = iota
	Enqueued
)

func (status EnqueuedStatus) String() string {
	return [...]string{"Failed", "Enqueued"}[status]
}
