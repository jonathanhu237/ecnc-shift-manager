package application

type contextKey string

const (
	requesterCtxKey contextKey = "requester"
	userCtxKey      contextKey = "user"
)
