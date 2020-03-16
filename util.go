package hfs

import (
	"log"
	"net/http"
	"strings"
)

func redirect(w http.ResponseWriter, newPath string) {
	w.Header().Set("Location", newPath)
	w.WriteHeader(http.StatusFound)
}

func forbidden(w http.ResponseWriter) {
	http.Error(w, "forbidden", http.StatusForbidden)
}

func notFound(w http.ResponseWriter) {
	http.Error(w, "not found", http.StatusNotFound)
}

func badRequest(w http.ResponseWriter) {
	http.Error(w, "bad request", http.StatusBadRequest)
}

func internalError(w http.ResponseWriter) {
	http.Error(w, "internal server error", http.StatusInternalServerError)
}

func stripPort(addr string) string {
	pos := strings.LastIndexByte(addr, ':')
	if pos > 0 {
		addr = addr[:pos]
	}
	return addr
}

func logf(logger *log.Logger, format string, v ...interface{}) {
	if logger != nil {
		logger.Printf(format, v...)
	}
}
