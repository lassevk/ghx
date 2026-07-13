package main

import "testing"

func TestParseOwner(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		// SCP-lignende SSH
		{"scp med .git", "git@github.com:lassevk/ghx.git", "lassevk", false},
		{"scp uten .git", "git@github.com:lassevk/ghx", "lassevk", false},
		{"scp org", "git@github.com:Larvik-Kommune/foo.git", "larvik-kommune", false},

		// SSH-URL
		{"ssh-url med .git", "ssh://git@github.com/lassevk/ghx.git", "lassevk", false},
		{"ssh-url uten .git", "ssh://git@github.com/lassevk/ghx", "lassevk", false},

		// HTTPS
		{"https med .git", "https://github.com/lassevk/ghx.git", "lassevk", false},
		{"https uten .git", "https://github.com/lassevk/ghx", "lassevk", false},
		{"https med port", "https://github.com:443/lassevk/ghx.git", "lassevk", false},

		// Case-insensitivitet
		{"blandet case owner", "https://github.com/LasseVK/Ghx.git", "lassevk", false},
		{"blandet case host", "git@GitHub.com:lassevk/ghx.git", "lassevk", false},

		// Ikke-github → feil
		{"gitlab https", "https://gitlab.com/lassevk/ghx.git", "", true},
		{"bitbucket scp", "git@bitbucket.org:lassevk/ghx.git", "", true},
		{"enterprise", "git@github.larvik.no:lassevk/ghx.git", "", true},

		// Ugyldig → feil
		{"tom", "", "", true},
		{"tull", "bare-noe-tekst", "", true},
		{"github uten owner", "https://github.com/", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseOwner(tt.url)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("parseOwner(%q) = %q, expected error", tt.url, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseOwner(%q) unexpected error: %v", tt.url, err)
			}
			if got != tt.want {
				t.Errorf("parseOwner(%q) = %q, expected %q", tt.url, got, tt.want)
			}
		})
	}
}
