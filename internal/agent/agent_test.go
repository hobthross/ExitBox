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

package agent

import (
	"os"
	"path/filepath"
	"testing"
)

func TestImportFile_Codex(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	srcFile := filepath.Join(srcDir, "config.toml")
	content := []byte("[model]\nprovider = \"custom\"\n")
	if err := os.WriteFile(srcFile, content, 0644); err != nil {
		t.Fatal(err)
	}

	c := &Codex{}
	if err := c.ImportFile(srcFile, dstDir); err != nil {
		t.Fatalf("ImportFile failed: %v", err)
	}

	target := filepath.Join(dstDir, ".codex", "config.toml")
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("expected file at %s: %v", target, err)
	}
	if string(data) != string(content) {
		t.Errorf("content mismatch: got %q, want %q", data, content)
	}
}

func TestImportFile_Claude(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	srcFile := filepath.Join(srcDir, "settings.json")
	content := []byte(`{"theme": "dark"}`)
	if err := os.WriteFile(srcFile, content, 0644); err != nil {
		t.Fatal(err)
	}

	c := &Claude{}
	if err := c.ImportFile(srcFile, dstDir); err != nil {
		t.Fatalf("ImportFile failed: %v", err)
	}

	target := filepath.Join(dstDir, ".claude", "settings.json")
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("expected file at %s: %v", target, err)
	}
	if string(data) != string(content) {
		t.Errorf("content mismatch: got %q, want %q", data, content)
	}
}

func TestImportFile_OpenCode(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	srcFile := filepath.Join(srcDir, "opencode.json")
	content := []byte(`{"provider": "anthropic"}`)
	if err := os.WriteFile(srcFile, content, 0644); err != nil {
		t.Fatal(err)
	}

	o := &OpenCode{}
	if err := o.ImportFile(srcFile, dstDir); err != nil {
		t.Fatalf("ImportFile failed: %v", err)
	}

	target := filepath.Join(dstDir, ".config", "opencode", "opencode.json")
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("expected file at %s: %v", target, err)
	}
	if string(data) != string(content) {
		t.Errorf("content mismatch: got %q, want %q", data, content)
	}
}

func TestConfigFilePath_Claude(t *testing.T) {
	c := &Claude{}
	got := c.ConfigFilePath("/ws/dir")
	want := filepath.Join("/ws/dir", ".claude", "settings.json")
	if got != want {
		t.Errorf("ConfigFilePath() = %q, want %q", got, want)
	}
}

func TestConfigFilePath_Codex(t *testing.T) {
	c := &Codex{}
	got := c.ConfigFilePath("/ws/dir")
	want := filepath.Join("/ws/dir", ".codex", "config.toml")
	if got != want {
		t.Errorf("ConfigFilePath() = %q, want %q", got, want)
	}
}

func TestConfigFilePath_OpenCode(t *testing.T) {
	o := &OpenCode{}
	got := o.ConfigFilePath("/ws/dir")
	want := filepath.Join("/ws/dir", ".config", "opencode", "opencode.json")
	if got != want {
		t.Errorf("ConfigFilePath() = %q, want %q", got, want)
	}
}

func TestImportFile_NonexistentSource(t *testing.T) {
	dstDir := t.TempDir()

	c := &Codex{}
	err := c.ImportFile("/nonexistent/file.toml", dstDir)
	if err == nil {
		t.Fatal("expected error for nonexistent source file")
	}
}
