package brokering

type QueueStatus int

const (
	Empty QueueStatus = iota
	Processing
)

func (status QueueStatus) String() string {
	return [...]string{"Empty", "Processing"}[status]
}
