package path_auth_redirector

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

type Config struct {
	Regex        string `json:"regex,omitempty"`
	Replacement  string `json:"replacement,omitempty"`
	HeaderName   string `json:"headerName,omitempty"`
	HeaderPrefix string `json:"headerPrefix,omitempty"`
}

func CreateConfig() *Config {
	return &Config{
		Regex:        "",
		Replacement:  "",
		HeaderName:   "",
		HeaderPrefix: "",
	}
}

type PathAuthRedirector struct {
	next   http.Handler
	config *Config
	regex  *regexp.Regexp
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	log.Printf("Initializing PathAuthRedirector with regex: %s, replacement: %s, headerName: %s, and headerPrefix: %s\n",
		config.Regex, config.Replacement, config.HeaderName, config.HeaderPrefix)
	regex, err := regexp.Compile(config.Regex)
	if err != nil {
		return nil, err
	}
	return &PathAuthRedirector{
		next:   next,
		config: config,
		regex:  regex,
	}, nil
}

func (p *PathAuthRedirector) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	matches := p.regex.FindStringSubmatch(req.URL.Path)
	if len(matches) > 1 {
		token := matches[1]
		log.Printf("Extracted token: %s\n", token)

		// Replace the matched part with the replacement
		newPath := p.regex.ReplaceAllString(req.URL.Path, p.config.Replacement)
		log.Printf("Modified request URL: %s\n", newPath)

		// Set the header with the extracted token, including the prefix if specified
		headerValue := token
		if p.config.HeaderPrefix != "" {
			headerValue = fmt.Sprintf("%s%s", p.config.HeaderPrefix, token)
		}
		req.Header.Set(p.config.HeaderName, headerValue)
		log.Printf("Set %s header: %s\n", p.config.HeaderName, headerValue)

		// Update the request URL
		req.URL.Path = newPath
		req.RequestURI = newPath
	}
	p.next.ServeHTTP(rw, req)
}
