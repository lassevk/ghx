package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// execGh kjører `gh` som en underprosess med GH_TOKEN satt til token og alle
// args uendret videre. stdin/stdout/stderr wires rett gjennom, og gh sin
// exit-kode propageres uendret.
//
// Dette erstatter den Unix-only syscall.Exec-sømmen fra #5 (se #8). exec.Command
// er portabel og gir Windows/macOS/Linux-paritet ut av boksen. Vi installerer
// bevisst ingen egen signalhåndtering: konsollen (Windows) og prosessgruppa
// (Unix) leverer Ctrl+C direkte til gh, som eier terminalen mens den kjører.
// Å legge et handler i ghx ville enten stjålet signalet fra gh eller doblet det.
//
// Ved suksess returnerer denne funksjonen aldri — den avslutter prosessen med
// gh sin exit-kode. Den returnerer bare en feil for tilfeller der gh ikke kunne
// startes i det hele tatt (ikke funnet i PATH, eller kunne ikke spawnes).
func execGh(args []string, token string) error {
	ghPath, err := exec.LookPath("gh")
	if err != nil {
		return fmt.Errorf("gh not found in PATH")
	}

	cmd := exec.Command(ghPath, args...)
	cmd.Env = buildEnv(token)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// gh kjørte, men avsluttet med ikke-null kode: speil koden uten å pakke
		// den inn i en «ghx:»-feilmelding, slik at wrapperen er transparent.
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}
		// gh kunne ikke startes (f.eks. binæren forsvant mellom LookPath og Run).
		return fmt.Errorf("could not run gh: %w", err)
	}

	os.Exit(0)
	return nil // ikke nåbar; tilfredsstiller kompilatoren
}
