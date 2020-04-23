package scraper

import (
	"fmt"
	"net/url"
	"strings"

	"go.uber.org/zap"
)

// Assumes all Path starts with a backslash and ends without one
type URL struct {
	Hostname string
	Path     string
}

func NewUrl(hostname, path string) *URL {
	if path == "" {
		path = "/"
	}
	return &URL{hostname, path}
}

func (u *URL) String() string {
	if u.Path == "/" {
		return fmt.Sprintf("http://%s/", u.Hostname)
	}
	return fmt.Sprintf("http://%s%s/", u.Hostname, u.Path)
}

func (u *URL) Concat() string {
	return fmt.Sprintf("%s%s", u.Hostname, u.Path)
}

var url_blacklist = []string{
	"javascript:void(0)",
}

func UrlParse(rawurl string) *URL {
	for _, blacklisted := range url_blacklist {
		if rawurl == blacklisted {
			return nil
		}
	}

	if strings.HasSuffix(rawurl, "/") {
		rawurl = rawurl[:len(rawurl)-1]
	}

	if !strings.Contains(rawurl, "//") {
		rawurl = "http://" + rawurl
	}

	parsedUrl, err := url.Parse(rawurl)
	if err != nil {
		zap.S().Debugf("net url.Parse error %v", err)
		return nil
	}

	if isMedia(parsedUrl) {
		zap.S().Debugf("is media, rawurl: %s", rawurl)
		return nil
	}

	if parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https" && parsedUrl.Scheme != "" {
		zap.S().Debugf("rawurl doesn't have proper scheme rawurl: %s", rawurl)
		return nil
	}

	if parsedUrl.Path == "" {
		parsedUrl.Path = "/"
	}

	u := URL{
		Hostname: parsedUrl.Hostname(),
		Path:     parsedUrl.Path,
	}

	if strings.HasPrefix(u.Hostname, "www.") {
		u.Hostname = u.Hostname[4:]
	}

	if parsedUrl.RawQuery != "" {
		u.Path = u.Path + "?" + parsedUrl.RawQuery
	}

	return &u
}

// Doesn't yet support dot notation, relative path
// Supports absolute path w/o domain name, and with domain
func ParseReferencedUrl(from *URL, rawurl string) *URL {
	parsedUrl := UrlParse(rawurl)
	if parsedUrl == nil {
		return nil
	}

	// Relative URL
	if parsedUrl.Hostname == "" {
		parsedUrl.Hostname = from.Hostname
	}
	return parsedUrl
}

func isMedia(u *url.URL) bool {
	media_extensions := []string{".png", ".mp3", ".mp4", ".ico", ".tif", ".jpeg", ".gif", ".jpg"}
	for _, ext := range media_extensions {
		if strings.HasSuffix(u.Path, ext) {
			return true
		}
	}
	return false
}
