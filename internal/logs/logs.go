package logs

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func RequestLogger(handler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.RawQuery)
		handler.ServeHTTP(w, r)
	})
}
