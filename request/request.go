package request

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type URL struct {
	ArgMap       map[string]string `yaml:"argument_map"`
	HeadersMap   map[string]string `yaml:"headers_map"`
	PayloadMap   map[string]string `yaml:"payload_map"`
	protocol     string            `yaml:"protocol"`
	domain       string            `yaml:"domain"`
	pathVals     []string          `yaml:"path_values"`
	path         string            `yaml:"path"`
	page         interface{}       `yaml:"page"`
	perPage      interface{}       `yaml:"per_page"`
	pageNoArg    string            `yaml:"page_no_arg"`
	perPageNoArg string            `yaml:"per_page_no_arg"`
	pagination   string            `yaml:"pagination"`
	full         string            `yaml:"full"`
	argStr       string            `yaml:"argStr"`
	method       string            `yaml:"method"`
	verbose      bool              `yaml:"verbose"`
	payload      io.Reader         `yaml:",omitempty"`
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
		verbose:      false,
	}
}

// VerboseOn turns on and off logging messages
func (u *URL) VerboseOn() *URL {
	u.verbose = true
	return u
}

// VerboseOff turns on and off logging messages
func (u *URL) VerboseOff() *URL {
	u.verbose = false
	return u
}

// Secure changes the http method to https
func (u *URL) Secure() *URL {
	u.protocol = "https://"
	if u.verbose {
		log.Println("https = on")
	}
	return u
}

// UnSecure changes the http method to http
func (u *URL) UnSecure() *URL {
	u.protocol = "http://"
	if u.verbose {
		log.Println("https = off")
	}
	return u
}

// SetPaginationArgNames changes the pagination name of the request
func (u *URL) SetPaginationArgNames(pageNo string, numberPerPage string) *URL {
	u.pageNoArg = pageNo
	u.perPageNoArg = numberPerPage
	if u.verbose {
		log.Println("https = off")
	}
	return u
}

// Pagination sets the pagination variables
func (u *URL) Pagination(page interface{}, per interface{}) *URL {
	u.page = page
	u.perPage = per
	u.pagination = fmt.Sprintf("%v=%v&%v=%v", u.pageNoArg, page, u.perPageNoArg, per)
	if u.verbose {
		log.Println("https = off")
	}
	return u
}

// Method sets the HTTP method
func (u *URL) Method(m string) *URL {
	u.method = m
	if u.verbose {
		log.Printf("method = %v", m)
	}
	return u
}

// Domain adds the domain to the request
func (u *URL) Domain(d string) *URL {
	u.domain = d
	if u.verbose {
		log.Printf("domain = %v", d)
	}
	return u
}

// Path accepts a number of strings and adjoins them by "/" in the order they are given and uses them as the path of a the request URL
func (u *URL) Path(pp ...string) *URL {
	path := ""
	u.pathVals = pp
	for _, p := range pp {
		path = fmt.Sprintf("%v/%v", path, p)
	}
	u.path = fmt.Sprint(u.path, path)
	if u.verbose {
		log.Printf("path = %v", u.path)
	}
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
	if u.verbose {
		log.Printf("query string args = %v", u.argStr)
	}
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
	if u.verbose {
		log.Println("Payload = ")
		for k, v := range u.PayloadMap {
			log.Printf("\t %v: %v\n", k, v)
		}
	}
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
	if u.verbose {
		log.Printf("full URL path  = %v\n", u.full)
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
