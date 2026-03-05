package agent

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyDir(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	// Create source structure
	_ = os.WriteFile(filepath.Join(src, "file1.txt"), []byte("hello"), 0644)
	_ = os.MkdirAll(filepath.Join(src, "subdir"), 0755)
	_ = os.WriteFile(filepath.Join(src, "subdir", "file2.txt"), []byte("world"), 0644)

	err := CopyDir(src, dst)
	if err != nil {
		t.Fatalf("CopyDir() error: %v", err)
	}

	// Verify
	data, err := os.ReadFile(filepath.Join(dst, "file1.txt"))
	if err != nil {
		t.Fatalf("reading file1.txt: %v", err)
	}
	if string(data) != "hello" {
		t.Errorf("file1.txt = %q, want %q", string(data), "hello")
	}

	data, err = os.ReadFile(filepath.Join(dst, "subdir", "file2.txt"))
	if err != nil {
		t.Fatalf("reading subdir/file2.txt: %v", err)
	}
	if string(data) != "world" {
		t.Errorf("subdir/file2.txt = %q, want %q", string(data), "world")
	}
}

func TestCopyDir_EmptyDir(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	err := CopyDir(src, dst)
	if err != nil {
		t.Fatalf("CopyDir(empty) error: %v", err)
	}
}

func TestCopyDir_NonexistentSrc(t *testing.T) {
	dst := t.TempDir()
	err := CopyDir("/nonexistent-path-xyz", dst)
	if err == nil {
		t.Error("CopyDir(nonexistent) should return error")
	}
}
