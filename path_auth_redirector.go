package path_auth_redirector

import (
	"context"
	"net/http"
	"regexp"
)

type Config struct {
	Regex           string `json:"regex,omitempty"`
	DefaultRedirect string `json:"defaultRedirect,omitempty"`
}

func CreateConfig() *Config {
	return &Config{
		Regex:           "",
		DefaultRedirect: "/",
	}
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		regex := regexp.MustCompile(config.Regex)
		matches := regex.FindStringSubmatch(req.URL.Path)
		if len(matches) > 1 {
			token := matches[1]
			req.Header.Set("Authorization", "Bearer "+token)
			next.ServeHTTP(rw, req)
		} else {
			http.Redirect(rw, req, config.DefaultRedirect, http.StatusFound)
		}
	}), nil
}
