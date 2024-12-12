package application

type contextKey string

type requester struct {
	id   int64
	role string
}

const (
	requesterCtxKey contextKey = "requester"
)
