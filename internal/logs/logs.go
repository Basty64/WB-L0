package logs

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func RequestLogger(handler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received requests: %s %s", r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
