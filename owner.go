package main

import (
	"fmt"
	"strings"
)

// parseOwnerRepo derives the GitHub owner (org or user) and the repo name from
// a git remote URL.
//
// It supports three forms:
//   - SCP-like SSH: git@github.com:Owner/Repo.git
//   - SSH URL:      ssh://git@github.com/Owner/Repo.git
//   - HTTPS:        https://github.com/Owner/Repo(.git)
//
// Only github.com is accepted as the host. Both owner and repo are normalized to
// lowercase since GitHub treats owner and repo names case-insensitively. A
// trailing ".git" suffix is stripped. Owner is always set on success; repo may
// be empty if the remote URL contains no repo name.
func parseOwnerRepo(remoteURL string) (owner, repo string, err error) {
	s := strings.TrimSpace(remoteURL)
	if s == "" {
		return "", "", fmt.Errorf("empty remote URL")
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
			return "", "", fmt.Errorf("could not parse remote URL: %s", s)
		}
		host, path = rest[:slash], rest[slash+1:]
		if colon := strings.Index(host, ":"); colon != -1 {
			host = host[:colon] // strip port
		}

	case strings.Contains(s, ":"):
		// SCP-like: [user@]host:path
		hostAndPath := s
		if at := strings.Index(s, "@"); at != -1 {
			hostAndPath = s[at+1:]
		}
		colon := strings.Index(hostAndPath, ":")
		host, path = hostAndPath[:colon], hostAndPath[colon+1:]

	default:
		return "", "", fmt.Errorf("could not parse remote URL: %s", s)
	}

	if strings.ToLower(host) != "github.com" {
		return "", "", fmt.Errorf("origin is not a github.com repository: %s", s)
	}

	path = strings.TrimPrefix(path, "/")
	ownerPart, repoPart, _ := strings.Cut(path, "/")
	ownerPart = strings.TrimSuffix(ownerPart, ".git")
	repoPart = strings.TrimSuffix(repoPart, ".git")
	if ownerPart == "" {
		return "", "", fmt.Errorf("could not determine owner from origin: %s", s)
	}
	return strings.ToLower(ownerPart), strings.ToLower(repoPart), nil
}
