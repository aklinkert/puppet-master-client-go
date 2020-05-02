package puppetmaster

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

type test struct {
	client   *Client
	apiToken string
	endpoint string
	server   *httptest.Server
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func newTestClient(t *testing.T, handler http.Handler) *test {
	apiToken := strconv.FormatUint(rand.Uint64(), 10)
	server := httptest.NewServer(handler)
	client, err := NewClient(server.URL, apiToken)
	if err != nil {
		t.Fatal("Failed to cvonstruct client:", err)
	}

	return &test{
		client:   client,
		apiToken: apiToken,
		server:   server,
		endpoint: server.URL,
	}
}

func dumbHandler(code int, body io.Reader) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(code)
		if body != nil {
			if _, err := io.Copy(rw, body); err != nil {
				panic(err)
			}
		}
	})
}

func readTestData(t *testing.T, filename string) []byte {
	filePath := fmt.Sprintf("test-data/%s", filename)
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read test data file %q: %v", filePath, err)
	}

	return b
}

func readJSONFileInto(t *testing.T, fileName string, data interface{}) {
	file := readTestData(t, fileName)
	if err := json.Unmarshal(file, data); err != nil {
		t.Fatalf("failed to unmarshal file into parameter: %v", err)
	}
}
