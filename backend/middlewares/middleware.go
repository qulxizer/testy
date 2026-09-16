package middlewares

import (
	"net/http"
	"slices"
)

type Middleware func(http.Handler) http.Handler

type Stack struct {
	middlewares []Middleware
}

func NewStack(ms ...Middleware) *Stack {
	return &Stack{middlewares: ms}
}

func (s *Stack) Then(finalHandler http.Handler) http.Handler {
	for _, middleware := range slices.Backward(s.middlewares) {
		finalHandler = middleware(finalHandler)
	}
	return finalHandler
}
