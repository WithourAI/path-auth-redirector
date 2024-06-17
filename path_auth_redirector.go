package path_auth_redirector

import (
	"context"
	"net/http"
	"regexp"
	"strings"
)

type Config struct {
	Regex    string `json:"regex,omitempty"`
	Redirect string `json:"redirect,omitempty"`
}

func CreateConfig() *Config {
	return &Config{
		Regex:    "",
		Redirect: "",
	}
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		regex := regexp.MustCompile(config.Regex)
		matches := regex.FindStringSubmatch(req.URL.Path)
		if len(matches) > 1 {
			token := matches[1]
			// Get the end position of the matched token in the URL path
			startIndex := strings.Index(req.URL.Path, token)
			endIndex := startIndex + len(token)
			remainingPath := req.URL.Path[endIndex:]
			req.Header.Set("Authorization", "Bearer "+token)
			req.URL.Path = config.Redirect + remainingPath
			next.ServeHTTP(rw, req)
		} else {
			// If there is no match, serve the request without redirection
			next.ServeHTTP(rw, req)
		}
	}), nil
}
