package feed

type CheckStatus struct {
	Status  Status
	Message string
}

type Status string

const (
	PendingStatus Status = ""
	SuccessStatus Status = "success"
	ErrorStatus   Status = "error"
)
