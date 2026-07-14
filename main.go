// Command ghx is a thin wrapper around `gh` that derives the repo's owner from
// the origin remote and runs `gh` with the right personal access token set in
// GH_TOKEN. See the README for background.
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "ghx: "+err.Error())
		os.Exit(1)
	}
}

// run performs the whole flow: derive the owner from origin, look up the token
// in the config, and exec `gh`. On any miss it returns an error and `gh` is
// never run.
func run(args []string) error {
	debug := os.Getenv("GHX_DEBUG") == "1"

	url, err := gitOriginURL()
	if err != nil {
		return err
	}
	if debug {
		fmt.Fprintf(os.Stderr, "ghx: origin = %s\n", url)
	}

	owner, repo, err := parseOwnerRepo(url)
	if err != nil {
		return err
	}
	repoKey := owner + "/" + repo
	if debug {
		if repo != "" {
			fmt.Fprintf(os.Stderr, "ghx: repo = %s\n", repoKey)
		} else {
			fmt.Fprintf(os.Stderr, "ghx: owner = %s (no repo in origin)\n", owner)
		}
	}

	tokens, err := loadConfig()
	if err != nil {
		return err
	}

	// Most specific wins: try owner/repo first, then fall back to owner.
	var token string
	switch {
	case repo != "" && tokens[repoKey] != "":
		token = tokens[repoKey]
		if debug {
			fmt.Fprintf(os.Stderr, "ghx: token found for '%s' (repo-specific)\n", repoKey)
		}
	case tokens[owner] != "":
		token = tokens[owner]
		if debug {
			label := "owner-level"
			if repo != "" {
				label = "owner-level fallback"
			}
			fmt.Fprintf(os.Stderr, "ghx: token found for '%s' (%s)\n", owner, label)
		}
	default:
		if repo != "" {
			return fmt.Errorf("no token configured for '%s' or '%s' in %s", repoKey, owner, configPath())
		}
		return fmt.Errorf("no token configured for '%s' in %s", owner, configPath())
	}

	// execGh only returns if something went wrong before gh took over the process.
	return execGh(args, token)
}

// buildEnv returns an environment like the current one, but with GH_TOKEN set to
// token and any existing GH_TOKEN/GITHUB_TOKEN removed, so there's no ambiguity
// about which token gh uses. Platform-neutral so step 2 (#8) can reuse it.
func buildEnv(token string) []string {
	base := os.Environ()
	env := make([]string, 0, len(base)+1)
	for _, e := range base {
		if strings.HasPrefix(e, "GH_TOKEN=") || strings.HasPrefix(e, "GITHUB_TOKEN=") {
			continue
		}
		env = append(env, e)
	}
	return append(env, "GH_TOKEN="+token)
}
