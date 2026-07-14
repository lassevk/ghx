package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// gitOriginURL returns the URL of the "origin" remote for the repo in the
// current directory. "Not in a git repo" and "no origin" are surfaced as
// distinct errors, since both are miss cases that should fail hard.
func gitOriginURL() (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		msg := stderr.String()
		switch {
		case strings.Contains(msg, "not a git repository"):
			return "", fmt.Errorf("not inside a git repository")
		case strings.Contains(strings.ToLower(msg), "no such remote"):
			return "", fmt.Errorf("repository has no 'origin' remote")
		default:
			return "", fmt.Errorf("git remote get-url origin failed: %v: %s", err, strings.TrimSpace(msg))
		}
	}
	return strings.TrimSpace(string(out)), nil
}
