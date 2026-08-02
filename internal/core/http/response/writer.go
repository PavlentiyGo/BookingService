package response

import "net/http"

type Writer struct {
	http.ResponseWriter
	StatusCode int
}

func (w *Writer) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
