package gorestdoc_test

import (
	apidoc "github.com/Holmes89/gorestdoc"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGenerateDoc(t *testing.T) {

	doc := apidoc.NewAPIDoc("Test", "This is a sample document for a test api")
	doc.AddDomain("Hello", "This is a test domain")
	doc.SetMarkdownFileName("test")

	helloHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("world!"))
	})

	router := mux.NewRouter()
	router.Handle("/hello", helloHandler).Methods("GET")
	ts := httptest.NewServer(router)
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL+"/hello", strings.NewReader("{\"hello\": \"world\"}"))

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
	//doc.GenerateHTMLFile()
	doc.GenerateMarkdownFile()
}
