package filesystem

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// LibexecPath computes the expected libexec path assuming a Filesystem
// Hierarchy Standard layout with the current executable located in the bin
// directory. It will return an error if the executable does not exist within
// the "bin" directory of such a layout, but it does not verify that the libexec
// directory exists.
func LibexecPath() (string, error) {
	// Compute the path to the current executable.
	executablePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("unable to compute executable path: %w", err)
	}

	if executablePath, err = filepath.EvalSymlinks(executablePath); err != nil {
		return "", fmt.Errorf("unable to resolve symlinks: %w", err)
	}

	// Check that the executable resides within a bin directory.
	if filepath.Base(filepath.Dir(executablePath)) != "bin" {
		return "", errors.New("executable does not reside within bin directory")
	}

	// Compute the expected libexec path.
	return filepath.Clean(filepath.Join(executablePath, "..", "..", "libexec")), nil
}
