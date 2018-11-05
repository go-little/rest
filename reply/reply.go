package reply

import (
	"encoding/json"
	"net/http"
)

// ErrorMaker interface
type ErrorMaker interface {
	Get(errorCode int, params ...string) (int, interface{})
}

// DefaultErrorMaker comment
var DefaultErrorMaker ErrorMaker

type reply struct {
	statusCode int
	headers    map[string]string
	body       []byte
	json       interface{}
}

type standardReply struct {
	reply
}

type errorReply struct {
	reply
}

// StatusCode comment
func StatusCode(statusCode int) *standardReply {
	return &standardReply{
		reply: reply{
			statusCode: statusCode,
			headers:    make(map[string]string),
			body:       nil,
			json:       nil,
		},
	}
}

// Error comment
func Error(errorCode int, params ...string) *errorReply {

	if DefaultErrorMaker == nil {
		panic("set DefaultErrorMaker for package github.com/go-little/rest/reply")
	}

	statusCode, json := DefaultErrorMaker.Get(errorCode, params...)

	return &errorReply{
		reply: reply{
			statusCode: statusCode,
			headers:    make(map[string]string),
			body:       nil,
			json:       json,
		},
	}
}

func (r *standardReply) Header(key string, value string) *standardReply {
	r.headers[key] = value
	return r
}

func (r *standardReply) Headers(headers map[string]string) *standardReply {
	r.headers = headers
	return r
}

func (r *standardReply) Body(body []byte) *standardReply {
	r.body = body
	r.json = nil
	return r
}

func (r *standardReply) JSON(json interface{}) *standardReply {
	r.json = json
	r.body = nil
	return r
}

func (r *reply) Do(w http.ResponseWriter) error {

	for key, value := range r.headers {
		w.Header().Set(key, value)
	}

	if r.json != nil {
		w.Header().Set("content-type", "application/json")
	}

	w.WriteHeader(r.statusCode)

	var err error
	if r.body != nil {
		_, err = w.Write(r.body)
		return err
	} else if r.json != nil {
		err = json.NewEncoder(w).Encode(r.json)
	} else {
		_, err = w.Write([]byte(""))
	}

	return err
}
