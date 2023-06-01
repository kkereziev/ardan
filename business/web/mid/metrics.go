package mid

import (
	"context"
	"net/http"

	"github.com/kkereziev/ardanlabs-service/business/sys/metrics"
	"github.com/kkereziev/ardanlabs-service/foundation/web"
)

func Metrics() web.Middleware {
	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// Add the metrics into the context for metric gathering.
			ctx = metrics.Set(ctx)

			// Handle updating the metrics that can be handled here.

			// Increment the request and goroutines counter.
			metrics.AddRequests(ctx)
			metrics.AddGoroutines(ctx)

			// Call the next handler.
			err := handler(ctx, w, r)
			if err != nil {
				// Increment if there is an error flowing through the request.
				metrics.AddErrors(ctx)
			}

			// Return the error so it can be handled further up the chain.
			return err
		}
	}
}
