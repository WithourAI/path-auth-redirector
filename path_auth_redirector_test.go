package path_auth_redirector

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestPathAuthRedirector(t *testing.T) {
	// Create a new instance of the plugin with the desired configuration
	config := &Config{
		Regex:        `/sk/(?P<token>[^/]+)(.*)`,
		Replacement:  "$2",
		HeaderName:   "Authorization",
		HeaderPrefix: "Bearer ",
	}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		log.Printf("Next handler called with path: %s\n", req.URL.Path)
	})

	handler, err := New(ctx, next, config, "path-auth-redirector")
	if err != nil {
		t.Fatalf("Failed to create plugin: %v", err)
	}

	testCases := []struct {
		name           string
		path           string
		expectedStatus int
		expectedHeader string
		expectedPath   string
	}{
		{
			name:           "Valid token in request path",
			path:           "/sk/validtoken123/resource",
			expectedStatus: http.StatusOK,
			expectedHeader: "Bearer validtoken123",
			expectedPath:   "/resource",
		},
		{
			name:           "Invalid request path",
			path:           "/invalid/path",
			expectedStatus: http.StatusOK,
			expectedHeader: "",
			expectedPath:   "/invalid/path",
		},
		{
			name:           "URL with any string token",
			path:           "/sk/anytoken456/endpoint",
			expectedStatus: http.StatusOK,
			expectedHeader: "Bearer anytoken456",
			expectedPath:   "/endpoint",
		},
		{
			name:           "URL with OpenAI-like token",
			path:           "/sk/sk-WHJajwidjldjjio289u90uaw/v1/chat/completions",
			expectedStatus: http.StatusOK,
			expectedHeader: "Bearer sk-WHJajwidjldjjio289u90uaw",
			expectedPath:   "/v1/chat/completions",
		},
		{
			name:           "URL with token containing special characters",
			path:           "/sk/" + url.PathEscape("sk_test_51AB-cD!ef@gh#ij$kl%mn^op") + "/v1/tokens",
			expectedStatus: http.StatusOK,
			expectedHeader: "Bearer sk_test_51AB-cD!ef@gh#ij$kl%mn^op",
			expectedPath:   "/v1/tokens",
		},
		{
			name:           "URL not starting with /sk",
			path:           "/batch/sk/sk-114514/v1/chat",
			expectedStatus: http.StatusOK,
			expectedHeader: "Bearer sk-114514",
			expectedPath:   "/batch/v1/chat",
		},
		{
			name:           "URL with token and query parameters",
			path:           "/sk/token789/api?param1=value1&param2=value2",
			expectedStatus: http.StatusOK,
			expectedHeader: "Bearer token789",
			expectedPath:   "/api",
		},
		{
			name:           "URL with multiple path segments after token",
			path:           "/sk/multitoken/segment1/segment2/segment3",
			expectedStatus: http.StatusOK,
			expectedHeader: "Bearer multitoken",
			expectedPath:   "/segment1/segment2/segment3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tc.path, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rw := httptest.NewRecorder()
			log.Printf("Running test case: %s\n", tc.name)

			handler.ServeHTTP(rw, req)

			if rw.Code != tc.expectedStatus {
				t.Errorf("Unexpected status code. Got %d, expected %d", rw.Code, tc.expectedStatus)
			}

			if req.Header.Get("Authorization") != tc.expectedHeader {
				t.Errorf("Unexpected Authorization header. Got %s, expected %s", req.Header.Get("Authorization"), tc.expectedHeader)
			}

			if req.URL.Path != tc.expectedPath {
				t.Errorf("Unexpected request path. Got %s, expected %s", req.URL.Path, tc.expectedPath)
			}

			log.Printf("Test case completed: %s\n", tc.name)
		})
	}
}
