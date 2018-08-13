package puppet_master

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

type test struct {
	client *Client
	apiKey string
	endpoint string
	server *httptest.Server
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func newTestClient(t *testing.T, handler http.Handler) *test {
	apiKey := strconv.FormatUint(rand.Uint64(), 10)
	server := httptest.NewServer(handler)
	client, err := NewClient(server.URL, apiKey)
	if err != nil {
		t.Fatal("Failed to cvonstruct client:", err)
	}

	return &test{
		client: client,
		apiKey: apiKey,
		server: server,
		endpoint: server.URL,
	}
}

func dumbHandler(code int, body io.Reader) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(code)
		if body != nil {
			io.Copy(rw, body)
		}
	})
}

func TestAddHeader(t *testing.T) {
	c := newTestClient(t, dumbHandler(200, nil))
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	t.Logf("Headers: %v", req.Header)
	c.client.addAuthentication(req)
	t.Logf("Headers: %v", req.Header)

	if len(req.Header) != 1 {
		t.Fatalf("Expected request to have 1 header, got %d", len(req.Header))
	}

	h := fmt.Sprintf("Bearer %v", c.apiKey)
	if req.Header.Get(authHeader) != h {
		t.Fatalf("Unexpected auth header %v, expected %v", req.Header.Get(authHeader), h)
	}
}

func TestBuildUrl(t *testing.T) {
	c := newTestClient(t, dumbHandler(200, nil))

	p := "jobs"
	q := map[string]string{
		"status": "test",
		"per_page": "100",
	}

	exp := fmt.Sprintf("%s/jobs?per_page=100&status=test", c.endpoint)
	res := c.client.buildUrl(p, q)

	t.Logf("Build url %v, expected %v", res, exp)
	if res != exp {
		t.Errorf("Expected to get url %v, got %v", exp, res)
	}
}

