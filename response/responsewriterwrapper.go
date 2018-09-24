package response

import "net/http"

type ResponseWriterWrapper struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

// NewResponseWriterWrapper comment
func NewResponseWriterWrapper(w http.ResponseWriter) *ResponseWriterWrapper {
	return &ResponseWriterWrapper{
		ResponseWriter: w,
		StatusCode:     http.StatusOK,
	}
}

func (r *ResponseWriterWrapper) WriteHeader(code int) {
	r.StatusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *ResponseWriterWrapper) Write(body []byte) (int, error) {
	r.Body = body
	return r.ResponseWriter.Write(body)
}
