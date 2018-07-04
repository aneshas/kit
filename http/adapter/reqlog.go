package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	gohttp "net/http"
	"time"

	"github.com/tonto/kit/http"
)

type logMessage struct {
	RemoteAddr string  `json:"remote_addr"`
	Method     string  `json:"method"`
	Path       string  `json:"path"`
	Took       string  `json:"took"`
	Body       *string `json:"body,omitempty"`
}

// WithRequestLogger creates a new request logging adapter
func WithRequestLogger(l *log.Logger, logRequestBody bool) http.Adapter {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(c context.Context, w gohttp.ResponseWriter, r *gohttp.Request) {
			msg := logMessage{
				RemoteAddr: r.RemoteAddr,
				Method:     r.Method,
				Path:       r.URL.Path,
			}

			if logRequestBody {
				switch r.Method {
				case gohttp.MethodPost, gohttp.MethodPut, gohttp.MethodPatch:
					buf, err := ioutil.ReadAll(r.Body)
					if err != nil {
						break
					}

					rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
					rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

					nbuf := new(bytes.Buffer)
					nbuf.ReadFrom(rdr1)
					body := nbuf.String()
					msg.Body = &body

					r.Body = rdr2
				}
			}

			defer func(t time.Time) {
				msg.Took = time.Since(t).String()
				data, _ := json.Marshal(&msg)
				l.Println(string(data))
			}(time.Now())

			h(c, w, r)
		}
	}
}
