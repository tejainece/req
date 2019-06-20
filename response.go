package req

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
)

type Response struct {
	*http.Response

	bodyParsed bool

	body []byte
}

func NewResponse(inner *http.Response) *Response {
	return &Response{Response: inner}
}

func (rs *Response) MediaType() string {
	val := rs.Header.Get("Content-Type")
	mt, _, err := mime.ParseMediaType(val)
	if err != nil {
		return ""
	}

	return mt
}

func (rs *Response) IsJson() bool {
	mt := rs.MediaType()

	return mt == "application/json"
}

func (rs *Response) IsXml() bool {
	mt := rs.MediaType()

	if mt == "text/xml" {
		return true
	}

	return mt == "application/xml"
}

func (rs *Response) Body() io.ReadCloser {
	bodyBytes, _ := rs.BodyBytes()
	return ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
}

func (rs *Response) BodyBytes() ([]byte, error) {
	if rs.bodyParsed {
		return rs.body, nil
	}

	b, err := ioutil.ReadAll(rs.Response.Body)
	rs.body = b
	rs.bodyParsed = true

	return rs.body, err
}

func (rs *Response) BodyString() (string, error) {
	b, err := rs.BodyBytes()
	if err != nil {
		return "", err
	}

	// TODO encoding

	return string(b), nil
}