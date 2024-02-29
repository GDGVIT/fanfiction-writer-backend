package main

import (
	"context"
	"net/http"

	"github.com/GDGVIT/fanfiction-writer-backend/fanfiction-backend/internal/data"
)

type contextKeyFF string

const userContextKey = contextKeyFF("user")

// contextSetUser() returns a new copy of the request with provided User struct
func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// contextGetUser() retrieves User struct from the request context
func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
