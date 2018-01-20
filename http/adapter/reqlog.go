package adapter

import (
	"context"
	"log"
	gohttp "net/http"
	"time"

	"github.com/tonto/kit/http"
)

const (
	gColor = "\x1b[32;1m"
	yColor = "\x1b[33;1m"
	bColor = "\x1b[34;1m"
	wColor = "\x1b[37;1m"
	nColor = "\x1b[0m"
)

// WithRequestLogger creates a new request logging adapter
func WithRequestLogger(l *log.Logger) http.Adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request) {
			t := time.Now()
			defer func() {
				l.Printf(
					"%s%s %s%s %s%s %s~%v%s",
					bColor, r.RemoteAddr,
					yColor, r.Method,
					gColor, r.URL.Path,
					wColor, time.Since(t),
					nColor,
				)
			}()
			h(c, w, r)
		}
	}
}
