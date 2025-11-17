package adapter

import (
	"net/url"
	"strings"
)

// normalizeDomainAndScheme normalizes a domain that may optionally include a scheme.
// - Returns domain without scheme or trailing slash
// - If scheme is provided, it overrides useHTTPS accordingly
// - If no scheme is provided, useHTTPS remains as defaultUseHTTPS
func normalizeDomainAndScheme(raw string, defaultUseHTTPS bool) (domain string, useHTTPS bool) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", defaultUseHTTPS
	}
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		if u, err := url.Parse(s); err == nil && u.Host != "" {
			d := strings.TrimSuffix(u.Host, "/")
			return d, strings.EqualFold(u.Scheme, "https")
		}
		// Fallback strip if parse fails
		s = strings.TrimPrefix(strings.TrimPrefix(s, "https://"), "http://")
	}
	return strings.TrimSuffix(s, "/"), defaultUseHTTPS
}

// encodePathSegments returns a URL-safe path by escaping each segment separately.
// It preserves slashes and does not add a leading slash.
func encodePathSegments(p string) string {
	trimmed := strings.TrimLeft(p, "/")
	if trimmed == "" {
		return ""
	}
	parts := strings.Split(trimmed, "/")
	for i := range parts {
		parts[i] = url.PathEscape(parts[i])
	}
	return strings.Join(parts, "/")
}
