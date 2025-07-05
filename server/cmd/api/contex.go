package main

import (
	"context"
	"biblioteca/internal/data"
	"net/http"
)


type contextKey string

const userContextKey = contextKey("user")

func (app *application) contextSetUser(r *http.Request, user *data.RegisterRequest) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *data.RegisterRequest {
	user, ok := r.Context().Value(userContextKey).(*data.RegisterRequest)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
