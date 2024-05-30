package publishing

type PublishStatus int

const (
	Failed PublishStatus = iota
	Published
)

func (status PublishStatus) String() string {
	return [...]string{"Failed", "Published"}[status]
}
