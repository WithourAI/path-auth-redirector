package path_auth_redirector

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMyPlugin(t *testing.T) {
	// Create a new instance of the plugin with the desired configuration
	config := &Config{
		Regex:    "^/sk/(?P<token>[^/]+).*",
		Redirect: "",
	}
	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})
	handler, err := New(ctx, next, config, "my-plugin")
	if err != nil {
		t.Fatalf("Failed to create plugin: %v", err)
	}

	// Test case 1: Valid token in the request path
	req1, _ := http.NewRequest("GET", "/sk/validtoken123/resource", nil)
	rw1 := httptest.NewRecorder()
	handler.ServeHTTP(rw1, req1)
	if rw1.Code != http.StatusOK {
		t.Errorf("Unexpected status code. Got %d, expected %d", rw1.Code, http.StatusOK)
	}
	expectedAuthHeader1 := "Bearer validtoken123"
	if req1.Header.Get("Authorization") != expectedAuthHeader1 {
		t.Errorf("Unexpected Authorization header. Got %s, expected %s", req1.Header.Get("Authorization"), expectedAuthHeader1)
	}
	expectedPath1 := "/resource"
	if req1.URL.Path != expectedPath1 {
		t.Errorf("Unexpected request path. Got %s, expected %s", req1.URL.Path, expectedPath1)
	}

	// Test case 2: Invalid request path
	req2, _ := http.NewRequest("GET", "/invalid/path", nil)
	rw2 := httptest.NewRecorder()
	handler.ServeHTTP(rw2, req2)
	if rw2.Code != http.StatusOK {
		t.Errorf("Unexpected status code. Got %d, expected %d", rw2.Code, http.StatusOK)
	}
	expectedRedirectURL2 := "/invalid/path"
	if req2.URL.Path != expectedRedirectURL2 {
		t.Errorf("Unexpected redirect URL. Got %s, expected %s", rw2.Header().Get("Location"), expectedRedirectURL2)
	}

	// Test case 3: URL with any string token
	req4, _ := http.NewRequest("GET", "/sk/anytoken456/endpoint", nil)
	rw4 := httptest.NewRecorder()
	handler.ServeHTTP(rw4, req4)
	if rw4.Code != http.StatusOK {
		t.Errorf("Unexpected status code. Got %d, expected %d", rw4.Code, http.StatusOK)
	}
	expectedAuthHeader4 := "Bearer anytoken456"
	if req4.Header.Get("Authorization") != expectedAuthHeader4 {
		t.Errorf("Unexpected Authorization header. Got %s, expected %s", req4.Header.Get("Authorization"), expectedAuthHeader4)
	}
	expectedPath4 := "/endpoint"
	if req4.URL.Path != expectedPath4 {
		t.Errorf("Unexpected request path. Got %s, expected %s", req4.URL.Path, expectedPath4)
	}
}
