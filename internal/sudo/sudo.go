// Package sudo provides functionalities to check if the current user can run sudo.
package sudo

import (
	"io"
	"os"

	"go.nhat.io/exec"
)

// Check checks if the current user can run sudo.
func Check() bool {
	_, err := exec.Run("sudo",
		exec.WithArgs("echo", "true"),
		exec.WithStdout(io.Discard),
		exec.WithStdin(os.Stdin),
		exec.WithStderr(os.Stderr),
	)

	return err == nil
}
