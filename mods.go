package req

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/pkg/errors"
)

type Mod interface {
	Before(req *Request) error
	After(resp *Response, req *Request) error
}

type Header struct {
	Key, value string
}

func (ma Header) Before(req *Request) error {
	req.Headers[ma.Key] = ma.value
	return nil
}

func (ma *Header) After(resp *Response, req *Request) error {
	return nil
}

type Headers map[string]string

func (ma Headers) Before(req *Request) error {
	for k, v := range ma {
		req.Headers[k] = v
	}
	return nil
}

func (ma *Headers) After(resp *Response, req *Request) error {
	return nil
}

type Decode struct {
	Value interface{}
}

func (ma Decode) Before(req *Request) error {
	return nil
}

func (ma *Decode) After(resp *Response, req *Request) error {
	respBytes, err := resp.BodyBytes()
	if err != nil {
		return errors.Wrapf(err, "error reading response body")
	}

	if resp.IsJson() {
		err = json.Unmarshal(respBytes, ma.Value)
		if err != nil {
			return errors.Wrapf(err, "error unmarshalling response")
		}
	} else if resp.IsXml() {
		err = xml.Unmarshal(respBytes, ma.Value)
		if err != nil {
			return errors.Wrapf(err, "error unmarshalling response")
		}
	} else {
		return fmt.Errorf("invalid content-type")
	}

	return nil
}

type DecodeOn200 struct {
	Value interface{}
}

func (ma DecodeOn200) Before(req *Request) error {
	return nil
}

func (ma DecodeOn200) After(resp *Response, req *Request) error {
	respBytes, err := resp.BodyBytes()
	if err != nil {
		return errors.Wrapf(err, "error reading response body")
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("statuscode %v body %v", resp.StatusCode, string(respBytes))
	}

	if resp.IsJson() {
		err = json.Unmarshal(respBytes, ma.Value)
		if err != nil {
			return errors.Wrapf(err, "error unmarshalling response")
		}
	} else if resp.IsXml() {
		err = xml.Unmarshal(respBytes, ma.Value)
		if err != nil {
			return errors.Wrapf(err, "error unmarshalling response")
		}
	} else {
		return fmt.Errorf("invalid content-type")
	}

	return nil
}

func Json(input interface{}) interface{} {
	bodyBytes, err := json.Marshal(input)
	if err != nil {
		return errors.Wrapf(err, "error marshalling")
	}

	return bodyBytes
}

func Xml(input interface{}) interface{} {
	bodyBytes, err := xml.Marshal(input)
	if err != nil {
		return errors.Wrapf(err, "error marshalling")
	}

	return bodyBytes
}

// MergeV merges multiple string maps into one
func MergeV(more ...map[string]string) (base map[string]string) {
	base = map[string]string{}

	for _, m := range more {
		for k, v := range m {
			base[k] = v
		}
	}

	return
}
