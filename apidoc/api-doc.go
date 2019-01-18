package apidoc

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// APIDoc Builds a base document for testing
type APIDoc struct {
	Title string
	Description string
	Generated time.Time
	Domains map[string]*Domain
}

// Domain fill in here
type Domain struct {
	Name string
	Description string
	Calls []*Call
}

// Call fill in here
type Call struct {
	Description string
	Method string
	RequestHeaders []string
	RequestBody string
	URL string
	RespCode int
	RespBody string
	RespHeaders []string
}

// emptyBody is an instance of empty reader.
var emptyBody = ioutil.NopCloser(strings.NewReader(""))

func NewAPIDoc(title, description string) *APIDoc {
	return &APIDoc{
		Title: title,
		Description: description,
		Domains: make(map[string]*Domain),
	}
}

// AddHTTPRequest
func (doc *APIDoc) AddHTTPRequest(domain, description string, req *http.Request) (*http.Response, error){

	d := doc.getDomain(domain)

	c := &Call{
		Description: description,
		Method: req.Method,
		URL: req.URL.Path,
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	save := resp.Body
	if resp.Body == nil {
		resp.Body = emptyBody
	} else {
		save, resp.Body, err = drainBody(resp.Body)
		if err != nil {
			return nil, err
		}
	}

	b, err := ioutil.ReadAll(save)
	if err != nil {
		return nil, err
	}

	c.RespCode = resp.StatusCode
	c.RespBody = string(b)

	d.Calls = append(d.Calls, c)

	return resp, nil
}

func (doc *APIDoc) AddDomain(domain, description string) {
	d := &Domain{
		Name: domain,
		Description: description,
	}
	doc.Domains[domain] = d
}

func (doc *APIDoc) getDomain(domain string) *Domain {
	d := doc.Domains[domain]
	if d == nil {
		d = &Domain{
			Name: domain,
		}
		doc.Domains[domain] = d
	}
	return d
}

func (doc *APIDoc) Print() {

	fmt.Printf("# %s\n", doc.Title)
	fmt.Printf("%s\n", doc.Description)

	for _, domain := range doc.Domains {

		fmt.Printf("## %s\n", domain.Name)
		fmt.Printf("%s\n", doc.Description)

		for _, call := range domain.Calls {

			fmt.Printf("### %s\n", call.Method)
			fmt.Printf("%s\n", call.Description)
			fmt.Println("*Request:*")
			fmt.Printf("`%s %s`\n", call.Method, call.URL)
			fmt.Println("*Response:*")
			fmt.Printf("```\nCode: %d\nBody: %s```", call.RespCode, call.RespBody)
		}
	}
}


// drainBody reads all of b to memory and then returns two equivalent
// ReadClosers yielding the same bytes.
//
// copied from httputil pkg
//
// It returns an error if the initial slurp of all bytes fails. It does not attempt
// to make the returned ReadClosers have identical error-matching behavior.
func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}