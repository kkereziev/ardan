package mid

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/kkereziev/ardanlabs-service/business/sys/metrics"
	"github.com/kkereziev/ardanlabs-service/foundation/web"
)

// Panics recovers from panics and converts the panic to an error so it is
// reported in Metrics and handled in Errors.
func Panics() web.Middleware {

	// This is the actual middleware function to be executed.
	return func(handler web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {

			// Defer a function to recover from a panic and set the err return
			// variable after the fact.
			defer func() {
				if rec := recover(); rec != nil {

					trace := debug.Stack()

					// Stack trace will be provided.
					err = fmt.Errorf("PANIC [%v] TRACE[%s]", rec, trace)

					// Updates the metrics stored in the context.
					metrics.AddPanic(ctx)
				}
			}()

			// Call the next handler and set its return value in the err variable.
			return handler(ctx, w, r)
		}
	}
}
