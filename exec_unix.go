//go:build unix

package main

import (
	"fmt"
	"os/exec"
	"syscall"
)

// execGh replaces the current process with `gh`, with GH_TOKEN set to token and
// all args passed through unchanged. On Unix we use syscall.Exec so that gh
// inherits stdin/stdout/stderr, tty, signals and exit code natively — the
// cleanest form of transparent passthrough.
//
// On success this function never returns (the process image is replaced).
// Step 2 (#8) swaps this seam to exec.Command for Windows support.
func execGh(args []string, token string) error {
	ghPath, err := exec.LookPath("gh")
	if err != nil {
		return fmt.Errorf("gh not found in PATH")
	}

	argv := append([]string{ghPath}, args...)
	return syscall.Exec(ghPath, argv, buildEnv(token))
}
