package feed

type Check struct {
	CheckStatus  Status
	CheckMessage string
}

type Status string

const (
	PendingStatus Status = ""
	SuccessStatus Status = "success"
	ErrorStatus   Status = "error"
)
