// Command ghx er en tynn wrapper rundt `gh` som utleder repoets owner fra
// origin-remoten og kjører `gh` med riktig personal access token satt i
// GH_TOKEN. Se README for bakgrunn.
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

// run utfører hele flyten: utled owner fra origin, slå opp token i config, og
// exec `gh`. Ved enhver «bom» returneres en feil og `gh` kjøres aldri.
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

	// Mest spesifikk vinner: prøv owner/repo først, fall så tilbake til owner.
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

	// execGh returnerer bare hvis noe gikk galt før gh overtok prosessen.
	return execGh(args, token)
}

// buildEnv returnerer et miljø likt det gjeldende, men med GH_TOKEN satt til
// token og eventuelle eksisterende GH_TOKEN/GITHUB_TOKEN fjernet, slik at det
// ikke er tvil om hvilket token gh bruker. Plattform-nøytral så steg 2 (#8)
// kan gjenbruke den.
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
