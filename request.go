package bounce

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type URL struct {
	ArgMap       map[string]string
	HeadersMap   map[string]string
	PayloadMap   map[string]string
	protocol     string
	domain       string
	path         string
	pageNoArg    string
	perPageNoArg string
	pagination   string
	full         string
	argStr       string
	method       string
	payload      io.Reader
}

func DefaultURL() *URL {
	return &URL{
		ArgMap:       make(map[string]string),
		HeadersMap:   make(map[string]string),
		PayloadMap:   make(map[string]string),
		protocol:     "http://",
		pageNoArg:    "page",
		perPageNoArg: "perPage",
		method:       "GET",
	}
}

// Secure changes the http method to https
func (u *URL) Secure() *URL {
	u.protocol = "https://"
	return u
}

// SetPaginationArgNames changes the pagination name of the request
func (u *URL) SetPaginationArgNames(pageNo string, numberPerPage string) *URL {
	u.pageNoArg = pageNo
	u.perPageNoArg = numberPerPage
	return u
}

// Pagination sets the pagination variables
func (u *URL) Pagination(page interface{}, per interface{}) *URL {
	u.pagination = fmt.Sprintf("%v=%v&%v=%v", u.pageNoArg, page, u.perPageNoArg, per)
	return u
}

// Method sets the HTTP method
func (u *URL) Method(m string) *URL {
	u.method = m
	return u
}

// Domain adds the domain to the request
func (u *URL) Domain(d string) *URL {
	u.domain = d
	return u
}

// Path accepts a number of strings and adjoins them by "/" in the order they are given and uses them as the path of a the request URL
func (u *URL) Path(pp ...string) *URL {
	path := ""
	for _, p := range pp {
		path = fmt.Sprintf("%v/%v", path, p)
	}
	u.path = fmt.Sprint(u.path, path)
	return u
}

// Args sets the query string arguments for a request
func (u *URL) Args(args map[string]string) *URL {
	var argStr string
	for k, v := range args {
		u.ArgMap[k] = v
	}
	for k, v := range u.ArgMap {
		argStr = fmt.Sprintf("%v%v=%v&", argStr, k, v)
	}
	u.argStr = strings.TrimSuffix(argStr, "&")
	return u
}

// Payload sets the json payload for the API
func (u *URL) Payload(pl map[string]string) *URL {
	PLStr := "{"
	for k, v := range pl {
		u.PayloadMap[k] = v
	}
	for k, v := range u.PayloadMap {
		PLStr = fmt.Sprintf("%v\"%v\":\"%v\",", PLStr, k, v)
	}
	PLStr = fmt.Sprint(strings.TrimSuffix(PLStr, ","), "}")
	u.payload = strings.NewReader(PLStr)
	return u
}

// Full drops out the appropriate pars for you to run http.NewRequest
func (u *URL) Full() (string, string, io.Reader) {
	u.full = fmt.Sprint(u.protocol, u.domain)
	if len(u.path) > 0 {
		u.full = fmt.Sprint(u.full, u.path)
	}
	sep := "?"
	if len(u.argStr) > 0 {
		u.full = fmt.Sprint(u.full, sep, u.argStr)
		sep = "&"
	}
	if len(u.pagination) > 0 {
		u.full = fmt.Sprint(u.full, sep, u.pagination)
	}
	return u.method, u.full, u.payload
}

// ConsumeAPI is the default method for firing the HTTP request
func (u *URL) ConsumeAPI(dest interface{}) error {

	req, err := http.NewRequest(u.Full())
	if err != nil {
		return nil
	}

	for k, v := range u.HeadersMap {
		req.Header.Add(k, v)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &dest)

	if err != nil {
		return err
	}
	return err
}
