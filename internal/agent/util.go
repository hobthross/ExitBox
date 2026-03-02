package agent

import (
	"fmt"
	"os"
	"path/filepath"
)

// CopyDir recursively copies directory contents from src to dst.
func CopyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	var errCount int
	var firstErr error
	for _, e := range entries {
		srcPath := filepath.Join(src, e.Name())
		dstPath := filepath.Join(dst, e.Name())
		if e.IsDir() {
			if err := os.MkdirAll(dstPath, 0755); err != nil {
				errCount++
				if firstErr == nil {
					firstErr = err
				}
				continue
			}
			if err := CopyDir(srcPath, dstPath); err != nil {
				errCount++
				if firstErr == nil {
					firstErr = err
				}
			}
		} else {
			data, err := os.ReadFile(srcPath)
			if err != nil {
				errCount++
				if firstErr == nil {
					firstErr = err
				}
				continue
			}
			if err := os.WriteFile(dstPath, data, 0644); err != nil {
				errCount++
				if firstErr == nil {
					firstErr = err
				}
			}
		}
	}
	if errCount > 0 {
		return fmt.Errorf("%d file(s) failed to copy, first error: %w", errCount, firstErr)
	}
	return nil
}
