package common

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
)

var AppBasePath = ""

func NormalizeAppBasePath(raw string) (string, error) {
	value := strings.TrimSpace(raw)
	if value == "" || value == "/" {
		return "", nil
	}
	if strings.ContainsAny(value, "?#\\") {
		return "", fmt.Errorf("must be a path without query, fragment, or backslash")
	}
	trimmed := strings.Trim(value, "/")
	if trimmed == "" {
		return "", nil
	}

	parts := strings.Split(trimmed, "/")
	cleanParts := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			continue
		}
		if part == "." || part == ".." {
			return "", fmt.Errorf("must not contain . or .. path segments")
		}
		cleanParts = append(cleanParts, part)
	}
	if len(cleanParts) == 0 {
		return "", nil
	}
	return "/" + strings.Join(cleanParts, "/"), nil
}

func AppBasePathCookiePath() string {
	if AppBasePath == "" {
		return "/"
	}
	return AppBasePath
}

func BuildAppPath(endpoint string) string {
	endpointPath, rawQuery, fragment := splitRelativeURL(endpoint)
	fullPath := joinURLPath(AppBasePath, endpointPath)
	if fullPath == "" {
		fullPath = "/"
	}
	if rawQuery != "" {
		fullPath += "?" + rawQuery
	}
	if fragment != "" {
		fullPath += "#" + fragment
	}
	return fullPath
}

func BuildPublicURL(serverAddress string, endpoint string) string {
	base := strings.TrimRight(strings.TrimSpace(serverAddress), "/")
	endpointPath, rawQuery, fragment := splitRelativeURL(endpoint)
	if base == "" {
		return BuildAppPath(endpoint)
	}

	parsed, err := url.Parse(base)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return buildFallbackPublicURL(base, endpointPath, rawQuery, fragment)
	}

	basePath := strings.TrimRight(parsed.Path, "/")
	if AppBasePath != "" && !pathHasSuffix(basePath, AppBasePath) {
		basePath = joinURLPath(basePath, AppBasePath)
	}
	parsed.Path = joinURLPath(basePath, endpointPath)
	parsed.RawPath = ""
	parsed.RawQuery = rawQuery
	parsed.Fragment = fragment
	return parsed.String()
}

func WithAppBasePath(next http.Handler) http.Handler {
	if AppBasePath == "" {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == AppBasePath {
			target := AppBasePath + "/"
			if r.URL.RawQuery != "" {
				target += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, target, http.StatusMovedPermanently)
			return
		}

		prefix := AppBasePath + "/"
		if !strings.HasPrefix(r.URL.Path, prefix) {
			http.NotFound(w, r)
			return
		}

		clone := new(http.Request)
		*clone = *r
		u := *r.URL
		u.Path = strings.TrimPrefix(r.URL.Path, AppBasePath)
		if u.Path == "" {
			u.Path = "/"
		}
		if r.URL.RawPath != "" {
			u.RawPath = strings.TrimPrefix(r.URL.RawPath, AppBasePath)
		}
		clone.URL = &u
		clone.RequestURI = u.RequestURI()
		next.ServeHTTP(w, clone)
	})
}

func splitRelativeURL(raw string) (string, string, string) {
	if raw == "" {
		return "", "", ""
	}
	parsed, err := url.Parse(raw)
	if err != nil || parsed.IsAbs() {
		return raw, "", ""
	}
	return parsed.Path, parsed.RawQuery, parsed.Fragment
}

func buildFallbackPublicURL(base string, endpointPath string, rawQuery string, fragment string) string {
	basePath := strings.TrimRight(base, "/")
	if AppBasePath != "" && !strings.HasSuffix(basePath, AppBasePath) {
		basePath += AppBasePath
	}
	result := strings.TrimRight(basePath, "/") + joinURLPath(endpointPath)
	if rawQuery != "" {
		result += "?" + rawQuery
	}
	if fragment != "" {
		result += "#" + fragment
	}
	return result
}

func pathHasSuffix(basePath string, suffix string) bool {
	if suffix == "" {
		return true
	}
	base := strings.TrimRight(basePath, "/")
	return base == suffix || strings.HasSuffix(base, suffix)
}

func joinURLPath(parts ...string) string {
	segments := make([]string, 0, len(parts))
	for _, part := range parts {
		for _, segment := range strings.Split(strings.Trim(part, "/"), "/") {
			if segment != "" {
				segments = append(segments, segment)
			}
		}
	}
	if len(segments) == 0 {
		return ""
	}
	return "/" + path.Join(segments...)
}
