//go:build unix

package main

import (
	"fmt"
	"os/exec"
	"syscall"
)

// execGh erstatter gjeldende prosess med `gh`, med GH_TOKEN satt til token og
// alle args uendret videre. På Unix bruker vi syscall.Exec slik at gh arver
// stdin/stdout/stderr, tty, signaler og exit-kode native — den reneste formen
// for transparent passthrough.
//
// Ved suksess returnerer denne funksjonen aldri (prosessbildet er byttet ut).
// Steg 2 (#8) bytter denne sømmen til exec.Command for Windows-støtte.
func execGh(args []string, token string) error {
	ghPath, err := exec.LookPath("gh")
	if err != nil {
		return fmt.Errorf("gh not found in PATH")
	}

	argv := append([]string{ghPath}, args...)
	return syscall.Exec(ghPath, argv, buildEnv(token))
}
