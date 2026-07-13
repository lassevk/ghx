# ghx

En tynn wrapper rundt [`gh`](https://cli.github.com/) som velger riktig
personal access token (PAT) automatisk ut fra hvilket repo du står i.

## Problemet

`gh` kan bytte mellom ulike GitHub-*kontoer*, men ikke mellom flere *tilganger
for samme konto*. Når man må bruke PAT-er (f.eks. fordi konto-wide innlogging
ikke er tillatt), og har flere owner-e (org-er/brukere) side om side, blir det
tungvint å jonglere hvilket token som gjelder hvor.

`ghx` løser det: den leser `origin`-remoten, finner owneren, slår opp riktig
PAT i en config-fil, setter den som `GH_TOKEN`, og kjører `gh` med alle
argumenter uendret.

## Bruk

```sh
ghx pr list
ghx api /user
ghx issue view 42
```

Alt sendes rett videre til `gh`, med `GH_TOKEN` satt til owner-ens token.
`ghx` krever å stå i et github.com-repo med en kjent owner — ellers feiler den
og kjører aldri `gh` (ingen risiko for feil konto).

Sett `GHX_DEBUG=1` for å se hvilken owner som ble utledet.

## Oppsett

1. Bygg: `go build -o ghx .` og legg binæren i PATH.
2. Opprett `~/.config/ghx/config.toml` (se `config.example.toml`):

   ```toml
   [owners]
   lassevk          = "ghp_..."
   "larvik-kommune" = "ghp_..."
   ```

3. Sikre fila: `chmod 600 ~/.config/ghx/config.toml`.

Nøkkelen er GitHub-owner (org/bruker) fra origin-URL-en; oppslag er
case-insensitivt. Respekterer `$XDG_CONFIG_HOME`.

## Status

MVP kjører på macOS/Linux via `syscall.Exec`. Windows-støtte kommer
([#8](https://github.com/lassevk/ghx/issues/8)).
