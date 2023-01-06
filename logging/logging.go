package logging

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
)

// Just for review info
// Logging package could be extended with dynamic logging
// which represent the option to change the level of logging (debug, warn, info, error)
// This help in generating fewer logs when not needed and set more logs when problem arise for debugging purposes
// This could be implemented through a configmap in the which is deployed in k8s cluster for examples separately from
// the application

// Package encapsulate productive json requirement logging which is required by a lot of
// analysing log tools

// The `key` will be available only in the log package, so no one will have access to such type
// therefore they can not change it in the context
type key string

const (
	requestIDKey    = key("requestId")
	RequestIdHeader = "X-Request-ID"
)

// Decorate adds uuid to request context so logs can be tracked.
// Used with Println to log the uuid with msg
func Decorate(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := uuid.New().String()
		ctx = context.WithValue(ctx, requestIDKey, id)
		f(w, r.WithContext(ctx))
	}
}

// DecorateHeader works as logging.Decorate but sets the uuid as RequestIdHeader
// so the client (consumer) could give unique problem id
func DecorateHeader(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := uuid.New().String()
		ctx = context.WithValue(ctx, requestIDKey, id)
		w.Header().Set(RequestIdHeader, id)
		f(w, r.WithContext(ctx))
	}
}

// Println prints message to os.Stdout with generated uuid and specific zerolog.Level
func Println(ctx context.Context, level zerolog.Level, msg string) {
	id, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		log.WithLevel(level).Msg(msg)
		log.Error().Msg("Could not find request id")
		return
	}
	fmt.Println()
	log.WithLevel(level).Str("requestId", id).Msg(msg)
}
