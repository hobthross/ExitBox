// ExitBox - Multi-Agent Container Sandbox
// Copyright (C) 2026 Cloud Exit B.V.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package fsutil

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// EnsureDir creates the directory (and parents) and returns its path.
func EnsureDir(parts ...string) string {
	p := filepath.Join(parts...)
	_ = os.MkdirAll(p, 0755)
	return p
}

// EnsureFile creates the parent dir and an empty JSON file if missing, returns its path.
func EnsureFile(parts ...string) string {
	p := filepath.Join(parts...)
	_ = os.MkdirAll(filepath.Dir(p), 0755)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		_ = os.WriteFile(p, []byte("{}\n"), 0644)
	}
	return p
}

// SeedDirOnce copies host dir into managed if host exists and managed is empty.
func SeedDirOnce(host, managed string) {
	if _, err := os.Stat(host); os.IsNotExist(err) {
		return
	}
	entries, err := os.ReadDir(managed)
	if err == nil && len(entries) > 0 {
		return
	}
	_ = os.MkdirAll(managed, 0755)
	_ = CopyDirRecursive(host, managed)
}

// SeedFileOnce copies host file to managed if host exists and managed does not.
func SeedFileOnce(host, managed string) {
	if _, err := os.Stat(host); os.IsNotExist(err) {
		return
	}
	if _, err := os.Stat(managed); err == nil {
		return
	}
	_ = os.MkdirAll(filepath.Dir(managed), 0755)
	data, err := os.ReadFile(host)
	if err == nil {
		_ = os.WriteFile(managed, data, 0644)
	}
}

// CopyDir recursively copies directory contents from src to dst.
// The dst directory must exist; subdirectories are created as needed.
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

// CopyDirRecursive copies a directory tree from src to dst.
// Symlinks are recreated rather than followed.
func CopyDirRecursive(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		if d.Type()&os.ModeSymlink != 0 {
			link, err := os.Readlink(path)
			if err != nil {
				return nil
			}
			_ = os.Remove(target)
			return os.Symlink(link, target)
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
