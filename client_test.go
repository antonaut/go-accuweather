package accuweather

import (
	"flag"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testKey = flag.String("test-key", "testtest", "accuweather test api key")

func TestMain(m *testing.M) {
	flag.Parse()
	m.Run()
}

// read test json data from testdata directory
// used to load test responses from the iex api
func readTestData(fileName string) (string, error) {
	b, err := ioutil.ReadFile("testdata/responses/" + fileName)
	if err != nil {
		return "", err
	}
	str := string(b)
	return str, nil
}

type mockHTTPClient struct {
	body    string
	headers map[string]string
	code    int
	err     error
}

func (c *mockHTTPClient) Get(url string) (*http.Response, error) {
	w := httptest.NewRecorder()
	w.WriteString(c.body)

	for key, value := range c.headers {
		w.Header().Add(key, value)
	}

	w.WriteHeader(c.code)

	resp := w.Result()
	return resp, c.err
}

func TestNewClient(t *testing.T) {
	c := NewClient(*testKey, &http.Client{})
	if c == nil {
		t.Fatalf("returned a nil client")
	}
}

func TestSearchLocation(t *testing.T) {
	body, err := readTestData("search_locations.json")
	if err != nil {
		t.Fatal(err)
	}

	httpc := mockHTTPClient{body: body, code: 200}
	c := NewClient(*testKey, &httpc)
	result, err := c.SearchForLocation("wichita")
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatalf("result was unexpectedly nil")
	}
}
