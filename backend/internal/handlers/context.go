package handlers

type contextKey string

const (
	requesterCtxKey contextKey = "requester"
	userCtxKey      contextKey = "user"
)
