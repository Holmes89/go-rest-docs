package go_rest_docs

import (
	"bytes"
	"fmt"
	"github.com/russross/blackfriday"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// APIDoc Builds a base document for testing
type APIDoc struct {
	Title string
	Description string
	Generated time.Time
	markdownFileName string
	htmlFileName string
	domains map[string]*domain
}

// domain fill in here
type domain struct {
	Name string
	Description string
	Calls []*call
}

// call fill in here
type call struct {
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
		markdownFileName: "README.md",
		htmlFileName: "api.html",
		domains: make(map[string]*domain),
	}
}

// SetMarkdownFileName overrides default output name to given name
func (doc *APIDoc) SetMarkdownFileName(name string) {
	doc.markdownFileName = name + ".md" // TODO check to see if it has extension
}

// SetHTMLFileName overrides default output name to given name
func (doc *APIDoc) SetHTMLFileName(name string) {
	doc.htmlFileName = name + ".html" // TODO check to see if it has extension
}

// AddHTTPRequest
func (doc *APIDoc) AddHTTPRequest(domain, description string, req *http.Request) (*http.Response, error){

	d := doc.getDomain(domain)

	c := &call{
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

// AddDomain allows you to define a domain and description
func (doc *APIDoc) AddDomain(name, description string) {
	d := &domain{
		Name: name,
		Description: description,
	}
	doc.domains[name] = d
}

func (doc *APIDoc) getDomain(name string) *domain {
	d := doc.domains[name]
	if d == nil {
		d = &domain{
			Name: name,
		}
		doc.domains[name] = d
	}
	return d
}

// Print outputs a string response
func (doc *APIDoc) Print() string{

	builder := MarkDownBuilder{}

	builder.H1(doc.Title).Body(doc.Description)

	for _, domain := range doc.domains {

		builder.H2(domain.Name).Body(domain.Description)

		for _, call := range domain.Calls {

			reqString := fmt.Sprintf("%s %s", call.Method, call.URL)
			respString := fmt.Sprintf("Code: %d\n\nBody: %s", call.RespCode, call.RespBody)
			builder.H3(call.Method).Body(call.Description).H4("Request").Code(reqString).H4("Response").Code(respString)
		}
	}

	return builder.Build()
}

// GenerateHTMLFile creates an HTML file from document struct
func (doc *APIDoc) GenerateHTMLFile() {
	md := doc.Print()
	output := blackfriday.Run([]byte(md), blackfriday.WithExtensions(blackfriday.CommonExtensions))

	f, err := os.Create(doc.htmlFileName)
	if err != nil {
		log.Fatal("could not create file")
	}
	defer f.Close()
	f.Write(output)
}

// GenerateHTMLFile creates an Markdown file from document struct
func (doc *APIDoc) GenerateMarkdownFile() {
	md := doc.Print()
	f, err := os.Create(doc.markdownFileName)
	if err != nil {
		log.Fatal("could not create file")
	}
	defer f.Close()
	f.Write([]byte(md))
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