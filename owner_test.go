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
		// SCP-lignende SSH
		{"scp med .git", "git@github.com:lassevk/ghx.git", "lassevk", "ghx", false},
		{"scp uten .git", "git@github.com:lassevk/ghx", "lassevk", "ghx", false},
		{"scp org", "git@github.com:Larvik-Kommune/foo.git", "larvik-kommune", "foo", false},

		// SSH-URL
		{"ssh-url med .git", "ssh://git@github.com/lassevk/ghx.git", "lassevk", "ghx", false},
		{"ssh-url uten .git", "ssh://git@github.com/lassevk/ghx", "lassevk", "ghx", false},

		// HTTPS
		{"https med .git", "https://github.com/lassevk/ghx.git", "lassevk", "ghx", false},
		{"https uten .git", "https://github.com/lassevk/ghx", "lassevk", "ghx", false},
		{"https med port", "https://github.com:443/lassevk/ghx.git", "lassevk", "ghx", false},

		// Case-insensitivitet — både owner og repo lowercases
		{"blandet case owner+repo", "https://github.com/LasseVK/Ghx.git", "lassevk", "ghx", false},
		{"blandet case host", "git@GitHub.com:lassevk/ghx.git", "lassevk", "ghx", false},

		// Owner uten repo → owner satt, repo tomt
		{"owner uten repo", "git@github.com:lassevk", "lassevk", "", false},

		// Ikke-github → feil
		{"gitlab https", "https://gitlab.com/lassevk/ghx.git", "", "", true},
		{"bitbucket scp", "git@bitbucket.org:lassevk/ghx.git", "", "", true},
		{"enterprise", "git@github.larvik.no:lassevk/ghx.git", "", "", true},

		// Ugyldig → feil
		{"tom", "", "", "", true},
		{"tull", "bare-noe-tekst", "", "", true},
		{"github uten owner", "https://github.com/", "", "", true},
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
