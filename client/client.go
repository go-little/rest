package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/go-little/rest/tracer"
)

type HTTPClient struct {
	ctx           context.Context
	name          string
	method        string
	url           string
	timeout       time.Duration
	retryAttempts int
	retryDelay    time.Duration
	retryRuleF    func(request *HTTPClient, response *HTTPResponse, err error) bool
	param         map[string]string
	query         url.Values
	header        http.Header
	form          url.Values
	body          []byte
	startTime     time.Time
}

type HTTPResponse struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

// NewHTTPClient func
func NewHTTPClient(ctx context.Context, name string) *HTTPClient {
	HTTPClient := &HTTPClient{
		ctx:           ctx,
		name:          name,
		method:        "GET",
		timeout:       2 * time.Second,
		retryAttempts: 0,
		param:         make(map[string]string),
		query:         url.Values{},
		header:        http.Header{},
		form:          url.Values{},
	}
	return HTTPClient
}

func add(current []string, values ...interface{}) []string {
	if current == nil {
		current = make([]string, 0)
	}

	for _, value := range values {
		current = append(current, fmt.Sprintf("%v", value))
	}
	return current
}

func (r *HTTPClient) Method(method string) *HTTPClient {
	r.method = method
	return r
}

func (r *HTTPClient) URL(url string) *HTTPClient {
	r.url = url
	return r
}

func (r *HTTPClient) Timeout(timeout time.Duration) *HTTPClient {
	r.timeout = timeout
	return r
}

func (r *HTTPClient) Retry(attempts int, delay time.Duration, ruleF func(request *HTTPClient, HTTPResponse *HTTPResponse, err error) bool) *HTTPClient {
	r.retryAttempts = attempts
	r.retryDelay = delay
	r.retryRuleF = ruleF
	return r
}

func (r *HTTPClient) Param(param map[string]string) *HTTPClient {
	r.param = param
	return r
}

func (r *HTTPClient) AddParam(name string, value interface{}) *HTTPClient {
	r.param[name] = fmt.Sprintf("%v", value)
	return r
}

func (r *HTTPClient) Query(query map[string][]string) *HTTPClient {
	r.query = query
	return r
}

func (r *HTTPClient) AddQuery(name string, value ...interface{}) *HTTPClient {
	r.query[name] = add(r.query[name], value...)
	return r
}

func (r *HTTPClient) Header(header map[string][]string) *HTTPClient {
	r.header = header
	return r
}

func (r *HTTPClient) AddHeader(name string, value ...interface{}) *HTTPClient {
	r.header[name] = add(r.header[name], value...)
	return r
}

func (r *HTTPClient) Form(form map[string][]string) *HTTPClient {
	r.form = form
	return r
}

func (r *HTTPClient) AddForm(name string, value ...interface{}) *HTTPClient {
	r.form[name] = add(r.form[name], value...)
	return r
}

func (r *HTTPClient) Body(body []byte) *HTTPClient {
	r.body = body
	return r
}

func (r *HTTPClient) JSONBody(body interface{}) *HTTPClient {
	r.AddHeader("content-type", "application/json")
	b, _ := json.Marshal(body)
	r.body = b
	return r
}

func (r *HTTPClient) Send() (*HTTPResponse, error) {
	r.startTime = time.Now()
	return r.send(r.retryAttempts)
}

func (r *HTTPClient) buildRequest() (*http.Request, error) {
	urlParsed, err := url.Parse(r.url)
	if err != nil {
		return nil, err
	}

	urlParsed.RawQuery = r.query.Encode()

	req, err := http.NewRequest(r.method, urlParsed.String(), bytes.NewReader(r.body))
	if err != nil {
		return nil, err
	}

	req.Header = r.header
	req.Form = r.form

	return req, nil
}

func (r *HTTPClient) send(attempts int) (*HTTPResponse, error) {
	var err error

	req, err := r.buildRequest()
	if err != nil {
		return nil, err
	}

	resSegment := &http.Response{}

	segment := tracer.StartExternalSegment(r.ctx, r.name, req)
	defer segment.End(resSegment)

	segment.Attr("request_header", req.Header)
	segment.Attr("request_method", req.Method)
	segment.Attr("request_url", fmt.Sprintf("%s://%s:%s%s", req.URL.Scheme, req.URL.Hostname(), req.URL.Port(), req.URL.Path))
	segment.Attr("request_querystring", req.URL.RawQuery)
	segment.Attr("request_body", string(r.body))
	segment.Attr("request_form", req.Form)

	var resp *HTTPResponse

	var body []byte
	var res *http.Response

	transport := http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.Dial(network, addr)
		},
	}

	HTTPClient := http.Client{
		Transport: &transport,
		Timeout:   r.timeout,
	}

	res, err = HTTPClient.Do(req)

	if err == nil {

		*resSegment = *res

		defer res.Body.Close()

		body, err = ioutil.ReadAll(res.Body)

		if err == nil {
			resp = &HTTPResponse{
				StatusCode: res.StatusCode,
				Header:     res.Header,
				Body:       body,
			}

			segment.Attr("response_header", resp.Header)
			segment.Attr("response_status_code", resp.StatusCode)
			segment.Attr("response_body", string(resp.Body))
		}
	}

	if attempts > 0 {
		if r.retryRuleF(r, resp, err) {
			fmt.Printf("Retry attempt %d\n", attempts)
			time.Sleep(r.retryDelay)
			return r.send(attempts - 1)
		}
	}

	segment.Attr("response_attempts", r.retryAttempts-attempts)
	segment.Attr("response_elapsed_milliseconds", time.Now().Sub(r.startTime)/time.Millisecond)

	return resp, err
}
