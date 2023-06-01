package testgrp

import (
	"context"
	"errors"
	"math/rand"
	"net/http"

	"github.com/kkereziev/ardanlabs-service/foundation/web"
	"go.uber.org/zap"
)

// Handlers manages the set of test endpoints.
type Handlers struct {
	Log *zap.SugaredLogger
}

// Test handler is for development.
func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if num := rand.Intn(100); num%2 == 0 {
		return errors.New("unknown error")
	}

	status := struct {
		Status string
	}{
		Status: "OK",
	}

	w.Header().Set("Content-Type", "application/json")

	h.Log.Infow("readiness", "statusCode", http.StatusOK, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)

	return web.Respond(ctx, w, status, http.StatusOK)
}
