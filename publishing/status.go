package publishing

type Status int

const (
	Failed Status = iota
	Published
)

func (status Status) String() string {
	return [...]string{"Failed", "Published"}[status]
}
