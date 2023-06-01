package mid

import (
	"context"
	"net/http"
	"time"

	"github.com/kkereziev/ardanlabs-service/foundation/web"
	"go.uber.org/zap"
)

func Logger(log *zap.SugaredLogger) web.Middleware {
	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// if the context is missing this value, request the service
			// to be shutdown gracefully
			v, err := web.GetValues(ctx)
			if err != nil {
				return err
			}

			log.Infow("request started", "traceid", v.TraceID, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)

			if err := handler(ctx, w, r); err != nil {
				return err
			}

			log.Infow("request completed", "traceid", v.TraceID, "method", r.Method,
				"path", r.URL.Path, "remoteaddr", r.RemoteAddr, "statuscode", v.StatusCode, "since", time.Since(v.Now))

			return nil
		}
	}
}
