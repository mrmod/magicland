package magicland

import (
	"encoding/json"
	"log"
	"net/http"
)

type httpErrorHandler func(string, error, http.ResponseWriter)

// HttpErrorHandlers
func badRequest(msg string, err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	writeError(msg, err, w)
}
func failedDependency(msg string, err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusFailedDependency)
	writeError(msg, err, w)
}
func writeError(msg string, err error, w http.ResponseWriter) {
	log.Printf("%s:%s", msg, err.Error())
	e := map[string]string{"error": err.Error()}
	b, _ := json.Marshal(e)
	w.Write(b)
}
