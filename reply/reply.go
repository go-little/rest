package reply

import (
	"encoding/json"
	"net/http"
)

type reply struct {
	statusCode int
	headers    map[string]string
	body       []byte
	json       interface{}
}

func StatusCode(statusCode int) *reply {
	return &reply{
		statusCode: statusCode,
		headers:    make(map[string]string),
		body:       nil,
		json:       nil,
	}
}

func (r *reply) Header(key string, value string) *reply {
	r.headers[key] = value
	return r
}

func (r *reply) Headers(headers map[string]string) *reply {
	r.headers = headers
	return r
}

func (r *reply) Body(body []byte) *reply {
	r.body = body
	r.json = nil
	return r
}

func (r *reply) JSON(json interface{}) *reply {
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
