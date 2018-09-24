package route

import (
	"net/http"
)

type Route struct {
	method      string
	pathPattern string
	handlerFunc func(w http.ResponseWriter, r *http.Request)
	auth        bool
}

type Routes []*Route

func (r *Routes) Add(method string, pathPattern string, handlerFunc func(w http.ResponseWriter, r *http.Request)) {
	r.add(method, pathPattern, handlerFunc, true)
}

func (r *Routes) AddWithAuth(method string, pathPattern string, handlerFunc func(w http.ResponseWriter, r *http.Request), auth bool) {
	r.add(method, pathPattern, handlerFunc, auth)
}

func (r *Routes) add(method string, pathPattern string, handlerFunc func(w http.ResponseWriter, r *http.Request), auth bool) {
	*r = append(*r, &Route{
		method:      method,
		pathPattern: pathPattern,
		handlerFunc: handlerFunc,
		auth:        auth,
	})
}
