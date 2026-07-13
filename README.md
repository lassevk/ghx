# ghx

A thin wrapper around [`gh`](https://cli.github.com/) that automatically picks
the right personal access token (PAT) based on which repository you're in.

## The problem

`gh` can switch between different GitHub *accounts*, but not between multiple
*accesses for the same account*. When you have to use PATs (e.g. because
account-wide login isn't permitted), and you have several owners (orgs/users)
side by side, juggling which token applies where becomes tedious.

`ghx` solves this: it reads the `origin` remote, determines the owner, looks up
the right PAT in a config file, sets it as `GH_TOKEN`, and runs `gh` with all
arguments passed through unchanged.

## Usage

```sh
ghx pr list
ghx api /user
ghx issue view 42
```

Everything is forwarded straight to `gh`, with `GH_TOKEN` set to the owner's
token. `ghx` requires being inside a github.com repository with a known owner —
otherwise it fails and never runs `gh` (no risk of using the wrong account).

Set `GHX_DEBUG=1` to see which owner was resolved.

## Setup

1. Build: `go build -o ghx .` and put the binary on your PATH.
2. Create `~/.config/ghx/config.toml` (see `config.example.toml`):

   ```toml
   [owners]
   lassevk          = "ghp_..."
   "larvik-kommune" = "ghp_..."
   ```

3. Secure the file: `chmod 600 ~/.config/ghx/config.toml`.

The key is the GitHub owner (org/user) from the origin URL; lookup is
case-insensitive. Respects `$XDG_CONFIG_HOME`.

## Status

Runs on Windows, macOS and Linux. `gh` is launched as a subprocess via
`exec.Command`, with stdin/stdout/stderr wired through and its exit code
propagated unchanged.
