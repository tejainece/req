package req

import (
	"bytes"
	"github.com/pkg/errors"
	"net/http"
)

type Request struct {
	method   string
	URL      *URL
	Headers  map[string]string
	Body     interface{}
	RespBody interface{}
	mods     []Mod
}

func (r *Request) Method() string {
	return r.method
}

func (r *Request) Do() (*Response, error) {
	for _, mod := range r.mods {
		if err := mod.Before(r); err != nil {
			return nil, err
		}
	}

	var bodyBytes []byte

	if r.Body == nil {
		// Ignore
	} else if data, ok := r.Body.([]byte); ok {
		bodyBytes = data
	} else if err, ok := r.Body.(error); ok {
		return nil, errors.Wrapf(err, "error in request body")
	} else if data, ok := r.Body.(string); ok {
		bodyBytes = []byte(data)
	} else {
		return nil, errors.Wrapf(err, "unsupported body type")
	}

	bodyReader := bytes.NewReader(bodyBytes)

	// TODO fmt.Println(string(bodyBytes))

	req, err := http.NewRequest(r.Method(), r.URL.String(), bodyReader)

	if err != nil {
		return nil, errors.Wrapf(err, "error initializing request")
	}

	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}

	var resp *Response
	{
		r, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		resp = NewResponse(r)
	}

	for i := range r.mods {
		mod := r.mods[len(r.mods)-i-1]
		if err := mod.After(resp, r); err != nil {
			return nil, err
		}
	}

	return resp, err
}
