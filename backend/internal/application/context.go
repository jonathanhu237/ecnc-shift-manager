package application

type contextKey string

type requester struct {
	id       int64
	username string
	role     string
	level    int
}

const (
	requesterCtxKey     contextKey = "requester"
	requesterDetailsKey contextKey = "requesterDetails"
)
