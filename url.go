package req

import (
	"fmt"
	"github.com/pkg/errors"
	"net/url"
)

type URL struct {
	url *url.URL
}

func NewURL() *URL {
	return &URL{}
}

func FromURI(uri string, query map[string]interface{}) (*URL, error) {
	ret := NewURL()
	err := ret.FromUri(uri, query)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (u URL) GetURL() (url.URL, bool) {
	if u.url == nil {
		return url.URL{}, false
	}
	return *u.url, true
}

func (u URL) String() string {
	if u.url == nil {
		return ""
	}

	return u.url.String()
}

func (u *URL) FromUri(uri string, query map[string]interface{}) error {
	parsed, err := url.ParseRequestURI(uri)
	if err != nil {
		return err
	}

	if query == nil {
		u.url = parsed
	}

	queryValue := url.Values{}

	for key, value := range query {
		var v string
		switch value.(type) {
		case string:
			v = value.(string)
		default:
			v = fmt.Sprintf("%v", value)
		}
		queryValue.Set(key, v)
	}
	parsed.RawQuery = queryValue.Encode()
	u.url = parsed

	return nil
}

func (u *URL) Get(headers map[string]string, mods ...Mod) (*Response, error) {
	if u.url == nil {
		return nil, errors.New("invalid URL")
	}

	r := &Request{
		method:  "GET",
		URL:     u,
		Headers: headers,
		mods:    mods,
	}

	return r.Do()
}

func (u *URL) Post(body interface{}, headers map[string]string, mods ...Mod) (*Response, error) {
	if u.url == nil {
		return nil, errors.New("invalid URL")
	}

	r := &Request{
		method:  "POST",
		URL:     u,
		Headers: headers,
		Body:    body,
		mods:    mods,
	}

	return r.Do()
}

func (u *URL) Put(body interface{}, headers map[string]string, mods ...Mod) (*Response, error) {
	if u.url == nil {
		return nil, errors.New("invalid URL")
	}

	r := &Request{
		method:  "PUT",
		URL:     u,
		Headers: headers,
		Body:    body,
		mods:    mods,
	}

	return r.Do()
}

func (u *URL) Delete(headers map[string]string, mods ...Mod) (*Response, error) {
	if u.url == nil {
		return nil, errors.New("invalid URL")
	}

	r := &Request{
		method:  "DELETE",
		URL:     u,
		Headers: headers,
		mods:    mods,
	}

	return r.Do()
}

func (u *URL) GetDecode(headers map[string]string, respBody interface{}, mods ...Mod) (*Response, error) {
	newMods := make([]Mod, len(mods))
	copy(newMods, mods)

	if respBody != nil {
		newMods = append(newMods, DecodeOn200{respBody})
	}

	return u.Get(headers, newMods...)
}

func (u *URL) PostDecode(body interface{}, headers map[string]string, respBody interface{}, mods ...Mod) (*Response, error) {
	newMods := make([]Mod, len(mods))
	copy(newMods, mods)

	if respBody != nil {
		newMods = append(newMods, DecodeOn200{respBody})
	}

	return u.Post(body, headers, newMods...)
}

func (u *URL) PutDecode(body interface{}, headers map[string]string, respBody interface{}, mods ...Mod) (*Response, error) {
	newMods := make([]Mod, len(mods))
	copy(newMods, mods)

	if respBody != nil {
		newMods = append(newMods, DecodeOn200{respBody})
	}

	return u.Put(body, headers, newMods...)
}

func (u *URL) DeleteDecode(headers map[string]string, respBody interface{}, mods ...Mod) (*Response, error) {
	newMods := make([]Mod, len(mods))
	copy(newMods, mods)

	if respBody != nil {
		newMods = append(newMods, DecodeOn200{respBody})
	}

	return u.Delete(headers, newMods...)
}