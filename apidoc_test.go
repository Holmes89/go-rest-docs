package main_test

import (
	"github.com/gorilla/mux"
	apidoc "go-rest-docs"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerateDoc(t *testing.T) {

	doc := apidoc.NewAPIDoc("Test", "This is a sample document for a test api")
	doc.AddDomain("Hello", "This is a test domain")

	helloHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("world!"))
	})

	router := mux.NewRouter()
	router.Handle("/hello", helloHandler).Methods("GET")
	ts := httptest.NewServer(router)
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL+"/hello", nil)

	desc := `This is supposed to be a long description with multiple lines and details that hopefully
will output properly and is representative of a longer description for what this will do`

	resp, err := doc.AddHTTPRequest("Hello", desc, req)

	if err != nil {
		t.Error("Error from test")
	}

	if resp.StatusCode != 200 {
		t.Errorf("Error code should be 200 not %d", resp.StatusCode)
	}

	t.Log(doc.Print())
	doc.GenerateHTMLFile()
	doc.GenerateMarkdownFile()
}
