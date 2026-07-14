package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// execGh runs `gh` as a subprocess with GH_TOKEN set to token and all args
// passed through unchanged. stdin/stdout/stderr are wired straight through, and
// gh's exit code is propagated unchanged.
//
// This replaces the Unix-only syscall.Exec seam from #5 (see #8). exec.Command
// is portable and gives Windows/macOS/Linux parity out of the box. We
// deliberately install no signal handling of our own: the console (Windows) and
// the process group (Unix) deliver Ctrl+C directly to gh, which owns the
// terminal while it runs. Adding a handler in ghx would either steal the signal
// from gh or double it.
//
// On success this function never returns — it exits the process with gh's exit
// code. It only returns an error for cases where gh could not be started at all
// (not found in PATH, or could not be spawned).
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
		// gh ran, but exited with a non-zero code: mirror the code without
		// wrapping it in a "ghx:" error message, so the wrapper stays transparent.
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}
		// gh could not be started (e.g. the binary vanished between LookPath and Run).
		return fmt.Errorf("could not run gh: %w", err)
	}

	os.Exit(0)
	return nil // unreachable; satisfies the compiler
}
