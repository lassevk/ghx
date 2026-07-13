package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// gitOriginURL returnerer URL-en til «origin»-remoten for repoet i gjeldende
// katalog. «Ikke i et git-repo» og «ingen origin» skilles ut som tydelige
// feil, siden begge er «bom»-tilfeller som skal føre til hard feil.
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
