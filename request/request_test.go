package request

import (
	"encoding/json"
	"regexp"
	"testing"
)

func Test_DefaultURL(t *testing.T) {
	d := DefaultURL()
	if d.protocol != "http://" {
		t.Errorf("Expected protocol to be \"http://\" got %v", d.protocol)
	}
	if d.pageNoArg != "page" {
		t.Errorf("Expected pageNoArg to be \"page\" got %v", d.pageNoArg)
	}
	if d.perPageNoArg != "perPage" {
		t.Errorf("Expected perPageNoArg to be \"perPage\" got %v", d.perPageNoArg)
	}
	if d.method != "GET" {
		t.Errorf("Expected method to be \"GET\" got %v", d.method)
	}
}

func Test_Secure(t *testing.T) {
	d := DefaultURL()
	if d.Secure().protocol != "https://" {
		t.Errorf("Expected protocol to be \"https\" got %v", d.protocol)
	}
}

func Test_SetPaginationArgNames(t *testing.T) {
	d := DefaultURL()
	d.SetPaginationArgNames("current", "amount")
	if d.pageNoArg != "current" {
		t.Errorf("Expected pageNoArg to be \"current\" got %v", d.pageNoArg)
	}
	if d.perPageNoArg != "amount" {
		t.Errorf("Expected perPageNoArg to be \"amount\" got %v", d.perPageNoArg)
	}
}

func Test_Pagination(t *testing.T) {
	d := DefaultURL()
	d.Pagination(63, 80)
	if d.pagination != "page=63&perPage=80" {
		t.Errorf("Expected pagination to be \"page=63&perPage=80\" got %v", d.pagination)
	}
}

func Test_Method(t *testing.T) {
	d := DefaultURL()
	if d.method != "GET" {
		t.Errorf("Expected method to be \"GET\" got %v", d.method)
	}
	d.Method("POST")
	if d.method != "POST" {
		t.Errorf("Expected method to be \"POST\" got %v", d.method)
	}
}

func Test_Domain(t *testing.T) {
	d := DefaultURL()
	if d.domain != "" {
		t.Errorf("Expected domain to be blank got %v", d.domain)
	}
	d.Domain("example.com")
	if d.domain != "example.com" {
		t.Errorf("Expected domain to be \"example.com\" got %v", d.domain)
	}
}

func Test_Path(t *testing.T) {
	d := DefaultURL()
	if d.path != "" {
		t.Errorf("Expected path to be blank got %v", d.path)
	}
	d.Path("here", "there", "everywhere")
	if d.path != "/here/there/everywhere" {
		t.Errorf("Expected path to be \"/here/there/everywhere\" got %v", d.path)
	}
}

func Test_Args(t *testing.T) {
	d := DefaultURL()
	if len(d.argStr) != 0 {
		t.Errorf("Expected args to be empty got %v", d.argStr)
	}
	m := make(map[string]string)
	m["what"] = "that"
	m["where"] = "there"
	m["when"] = "then"
	testStrings := [6]string{
		"where=there&when=then&what=that",
		"where=there&what=that&when=then",
		"when=then&where=there&what=that",
		"when=then&what=that&where=there",
		"what=that&when=then&where=there",
		"what=that&where=there&when=then",
	}
	ok := false
	d.Args(m)
	for _, v := range testStrings {
		if d.argStr == v {
			ok = true
		}
	}
	if !ok {
		t.Errorf("Expected args to be like \"what=that&where=there&when=then\" got %v", d.argStr)
	}
}

func Test_Payload(t *testing.T) {
	type Thing struct {
		What, Where, When string
	}
	d := DefaultURL()
	m := make(map[string]string)
	m["what"] = "that"
	m["where"] = "there"
	m["when"] = "then"
	d.Payload(m)
	dec := json.NewDecoder(d.payload)
	var j Thing
	err := dec.Decode(&j)
	if err != nil {
		t.Errorf("decoding error Test_%v", "Payload")
	}
	if j.When != m["when"] && j.Where != m["where"] && j.What != m["what"] {
		t.Errorf("Expected args to be like { that there then }, got %v", j)
	}
}

func Test_Full(t *testing.T) {
	d := DefaultURL()
	m := make(map[string]string)
	m["what"] = "that"
	m["where"] = "there"
	m["when"] = "then"
	_, f, _ := d.Secure().Domain("example.com").Path("here", "there", "everywhere").Args(m).Pagination(1, 100).Full()

	re := regexp.MustCompile(`https:\/\/example.com\/here\/there\/everywhere\?\D*=\D*&\D*=\D*&\D*=\D*=\d&\D*=\d{3}`)

	if len(re.FindStringIndex(f)) == 0 {
		t.Errorf("Expected full path to be like \"https://example.com/here/there/everywhere?what=that&where=there&when=then&page=1&perPage=100\" got %v", d.full)
	}
}
