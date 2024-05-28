package subscribing

type Status int

const (
	Failed Status = iota
	Received
)

func (status Status) String() string {
	return [...]string{"Failed", "Received"}[status]
}
