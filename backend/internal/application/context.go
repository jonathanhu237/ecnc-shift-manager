package application

type contextKey string

type requester struct {
	id    int64
	role  string
	level int
}

const (
	requesterCtxKey contextKey = "requester"
)
