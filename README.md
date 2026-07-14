# ghx

A thin wrapper around [`gh`](https://cli.github.com/) that automatically picks
the right personal access token (PAT) based on which repository you're in.

## The problem

`gh` can switch between different GitHub *accounts*, but not between multiple
*accesses for the same account*. When you have to use PATs (e.g. because
account-wide login isn't permitted), and you have several owners (orgs/users)
side by side, juggling which token applies where becomes tedious.

`ghx` solves this: it reads the `origin` remote, determines the owner (and
repo), looks up the right PAT in a config file, sets it as `GH_TOKEN`, and runs
`gh` with all arguments passed through unchanged.

## Usage

```sh
ghx pr list
ghx api /user
ghx issue view 42
```

Everything is forwarded straight to `gh`, with `GH_TOKEN` set to the matching
token. `ghx` requires being inside a github.com repository with a known owner —
otherwise it fails and never runs `gh` (no risk of using the wrong account).

Set `GHX_DEBUG=1` to see which owner/repo was resolved and which token matched.

## Setup

1. Build: `go build -o ghx .` and put the binary on your PATH.
2. Create `~/.config/ghx/config.toml` (see `config.example.toml`):

   ```toml
   [tokens]
   lassevk          = "ghp_..."
   "larvik-kommune" = "ghp_..."
   "lassevk/ghx"    = "ghp_..."   # per-repo override
   ```

3. Secure the file: `chmod 600 ~/.config/ghx/config.toml`.

Keys map to PATs and lookup is case-insensitive. A key without a slash is an
owner (org/user); a key with a slash is a specific `owner/repo`. For a given
repo, ghx tries the `owner/repo` key first and falls back to the `owner` key —
the most specific match wins and replaces the owner token, so a single repo can
be given a more-privileged PAT without widening access to the whole owner.
Respects `$XDG_CONFIG_HOME`.

## Status

MVP runs on macOS/Linux via `syscall.Exec`. Windows support is pending
([#8](https://github.com/lassevk/ghx/issues/8)).
