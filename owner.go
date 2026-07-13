package main

import (
	"fmt"
	"strings"
)

// parseOwner utleder GitHub-owneren (org eller bruker) fra en git remote-URL.
//
// Støtter tre former:
//   - SCP-lignende SSH: git@github.com:Owner/Repo.git
//   - SSH-URL:          ssh://git@github.com/Owner/Repo.git
//   - HTTPS:            https://github.com/Owner/Repo(.git)
//
// Kun github.com aksepteres som host. Owner normaliseres til lowercase siden
// GitHub behandler owner-navn case-insensitivt. En eventuell «.git»-endelse
// strippes.
func parseOwner(remoteURL string) (string, error) {
	s := strings.TrimSpace(remoteURL)
	if s == "" {
		return "", fmt.Errorf("tom remote-URL")
	}

	var host, path string
	switch {
	case strings.Contains(s, "://"):
		// scheme://[user@]host[:port]/path
		rest := s[strings.Index(s, "://")+3:]
		if at := strings.Index(rest, "@"); at != -1 {
			rest = rest[at+1:]
		}
		slash := strings.Index(rest, "/")
		if slash == -1 {
			return "", fmt.Errorf("klarte ikke tolke remote-URL: %s", s)
		}
		host, path = rest[:slash], rest[slash+1:]
		if colon := strings.Index(host, ":"); colon != -1 {
			host = host[:colon] // strip port
		}

	case strings.Contains(s, ":"):
		// SCP-lignende: [user@]host:path
		hostAndPath := s
		if at := strings.Index(s, "@"); at != -1 {
			hostAndPath = s[at+1:]
		}
		colon := strings.Index(hostAndPath, ":")
		host, path = hostAndPath[:colon], hostAndPath[colon+1:]

	default:
		return "", fmt.Errorf("klarte ikke tolke remote-URL: %s", s)
	}

	if strings.ToLower(host) != "github.com" {
		return "", fmt.Errorf("origin er ikke et github.com-repo: %s", s)
	}

	path = strings.TrimPrefix(path, "/")
	owner, _, _ := strings.Cut(path, "/")
	owner = strings.TrimSuffix(owner, ".git")
	if owner == "" {
		return "", fmt.Errorf("klarte ikke utlede owner fra origin: %s", s)
	}
	return strings.ToLower(owner), nil
}
