package path_auth_redirector

import (
	"context"
	"log"
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

type PathAuthRedirector struct {
	next   http.Handler
	config *Config
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	log.Printf("Initializing PathAuthRedirector with regex: %s and redirect: %s\n", config.Regex, config.Redirect)
	return &PathAuthRedirector{
		next:   next,
		config: config,
	}, nil
}

func (p *PathAuthRedirector) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	regex := regexp.MustCompile(p.config.Regex)
	matches := regex.FindStringSubmatch(req.URL.Path)
	if len(matches) > 1 {
		token := matches[1]
		log.Printf("Extracted token: %s\n", token)

		// Get the end position of the matched token in the URL path
		startIndex := strings.Index(req.URL.Path, token)
		endIndex := startIndex + len(token)
		remainingPath := req.URL.Path[endIndex:]
		log.Printf("Remaining path after token: %s\n", remainingPath)

		req.Header.Set("Authorization", "Bearer "+token)
		req.URL.Path = p.config.Redirect + remainingPath
		req.RequestURI = req.URL.Path
		log.Printf("Modified request URL: %s\n", req.URL.Path)
		log.Printf("Set Authorization header: Bearer %s\n", token)

		p.next.ServeHTTP(rw, req)
	} else {
		p.next.ServeHTTP(rw, req)
	}
}
