package gorestdoc_test

import (
	"bytes"
	apidoc "github.com/Holmes89/gorestdoc"
	"github.com/gorilla/mux"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestGenerateDoc(t *testing.T) {

	doc := apidoc.NewAPIDoc("Test", "This is a sample document for a test api")
	doc.AddDomain("Hello", "This is a test domain")
	doc.SetMarkdownFileName("test")


	router := newTestServer()
	ts := httptest.NewServer(router)
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL+"/hello?queryParam=test", strings.NewReader("{\"hello\": \"world\"}"))

	desc := `This is supposed to be a long description with multiple lines and details that hopefully
will output properly and is representative of a longer description for what this will do`

	resp, err := doc.AddHTTPRequest("Hello", desc, req)

	// Nil body
	req, _ = http.NewRequest("GET", ts.URL+"/hello", nil)

	resp, err = doc.AddHTTPRequest("Hello", desc, req)

	if err != nil {
		t.Error("Error from test")
	}

	if resp.StatusCode != 200 {
		t.Errorf("Error code should be 200 not %d", resp.StatusCode)
	}


	// Test Form
	var formBytes bytes.Buffer
	w := multipart.NewWriter(&formBytes)

	f, err := os.Open("test.md")
	if err != nil {
		t.Error("unable to open file")
		return
	}

	fw, err := w.CreateFormFile("file", f.Name())
	if err != nil {
		t.Errorf("unable to create file")
		return
	}

	io.Copy(fw, f)

	w.WriteField("name", "test")
	w.WriteField("type", "type")

	w.Close()

	req, err = http.NewRequest("POST",  ts.URL+"/upload", &formBytes)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err = doc.AddHTTPRequest("Form", desc, req)

	if err != nil {
		log.Print(err.Error())
		t.Error("Error from test")
		return
	}

	if resp.StatusCode != 201 {
		t.Errorf("Error code should be 201 not %d", resp.StatusCode)
		return
	}

	doc.GenerateMarkdownFile()
}

func newTestServer() *mux.Router {

	helloHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("world!"))
	})

	postFormHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("created!"))
	})

	router := mux.NewRouter()
	router.Handle("/hello", helloHandler).Methods("GET")
	router.Handle("/upload", postFormHandler).Methods("POST")

	return router
}