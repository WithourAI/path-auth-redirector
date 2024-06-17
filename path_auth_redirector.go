package path_auth_redirector

import (
	"context"
	"net/http"
	"regexp"
)

type Config struct {
	Regex    string `json:"regex,omitempty"`
	Redirect string `json:"redirect,omitempty"`
}

func CreateConfig() *Config {
	return &Config{
		Regex:    "",
		Redirect: "/",
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
			http.Redirect(rw, req, config.Redirect, http.StatusFound)
		}
	}), nil
}
