package main

import "testing"

func TestParseOwnerRepo(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantOwner string
		wantRepo  string
		wantErr   bool
	}{
		// SCP-like SSH
		{"scp with .git", "git@github.com:lassevk/ghx.git", "lassevk", "ghx", false},
		{"scp without .git", "git@github.com:lassevk/ghx", "lassevk", "ghx", false},
		{"scp org", "git@github.com:Larvik-Kommune/foo.git", "larvik-kommune", "foo", false},

		// SSH URL
		{"ssh-url with .git", "ssh://git@github.com/lassevk/ghx.git", "lassevk", "ghx", false},
		{"ssh-url without .git", "ssh://git@github.com/lassevk/ghx", "lassevk", "ghx", false},

		// HTTPS
		{"https with .git", "https://github.com/lassevk/ghx.git", "lassevk", "ghx", false},
		{"https without .git", "https://github.com/lassevk/ghx", "lassevk", "ghx", false},
		{"https with port", "https://github.com:443/lassevk/ghx.git", "lassevk", "ghx", false},

		// Case-insensitivity — both owner and repo are lowercased
		{"mixed case owner+repo", "https://github.com/LasseVK/Ghx.git", "lassevk", "ghx", false},
		{"mixed case host", "git@GitHub.com:lassevk/ghx.git", "lassevk", "ghx", false},

		// Owner without repo → owner set, repo empty
		{"owner without repo", "git@github.com:lassevk", "lassevk", "", false},

		// Non-github → error
		{"gitlab https", "https://gitlab.com/lassevk/ghx.git", "", "", true},
		{"bitbucket scp", "git@bitbucket.org:lassevk/ghx.git", "", "", true},
		{"enterprise", "git@github.larvik.no:lassevk/ghx.git", "", "", true},

		// Invalid → error
		{"empty", "", "", "", true},
		{"garbage", "just-some-text", "", "", true},
		{"github without owner", "https://github.com/", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, err := parseOwnerRepo(tt.url)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("parseOwnerRepo(%q) = (%q, %q), expected error", tt.url, owner, repo)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseOwnerRepo(%q) unexpected error: %v", tt.url, err)
			}
			if owner != tt.wantOwner || repo != tt.wantRepo {
				t.Errorf("parseOwnerRepo(%q) = (%q, %q), expected (%q, %q)", tt.url, owner, repo, tt.wantOwner, tt.wantRepo)
			}
		})
	}
}
