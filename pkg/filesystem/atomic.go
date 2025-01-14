package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	// atomicWriteTemporaryNamePrefix is the file name prefix to use for
	// intermediate temporary files used in atomic writes.
	atomicWriteTemporaryNamePrefix = TemporaryNamePrefix + "atomic-write"
)

// WriteFileAtomic writes a file to disk in an atomic fashion by using an
// intermediate temporary file that is swapped in place using a rename
// operation.
func WriteFileAtomic(path string, data []byte, permissions os.FileMode) error {
	// Create a temporary file. The os package already uses secure permissions
	// for creating temporary files, so we don't need to change them.
	temporary, err := os.CreateTemp(filepath.Dir(path), atomicWriteTemporaryNamePrefix)
	if err != nil {
		return fmt.Errorf("unable to create temporary file: %w", err)
	}

	// Write data.
	if _, err = temporary.Write(data); err != nil {
		temporary.Close()
		os.Remove(temporary.Name())
		return fmt.Errorf("unable to write data to temporary file: %w", err)
	}

	// Close out the file.
	if err = temporary.Close(); err != nil {
		os.Remove(temporary.Name())
		return fmt.Errorf("unable to close temporary file: %w", err)
	}

	// Set the file's permissions.
	if err = os.Chmod(temporary.Name(), permissions); err != nil {
		os.Remove(temporary.Name())
		return fmt.Errorf("unable to change file permissions: %w", err)
	}

	// Rename the file.
	if err = os.Rename(temporary.Name(), path); err != nil {
		os.Remove(temporary.Name())
		return fmt.Errorf("unable to rename file: %w", err)
	}

	// Success.
	return nil
}
