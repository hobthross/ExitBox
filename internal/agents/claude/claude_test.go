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

package claude

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestClaudeAgent(t *testing.T) {
	c := &Claude{}

	if c.Name() != "claude" {
		t.Errorf("Name() = %q, want %q", c.Name(), "claude")
	}
	if c.DisplayName() != "Claude Code" {
		t.Errorf("DisplayName() = %q, want %q", c.DisplayName(), "Claude Code")
	}

	// HostConfigPaths
	paths := c.HostConfigPaths()
	if len(paths) != 1 {
		t.Fatalf("HostConfigPaths() returned %d paths, want 1", len(paths))
	}
	if !strings.HasSuffix(paths[0], ".claude") {
		t.Errorf("HostConfigPaths()[0] = %q, should end with .claude", paths[0])
	}

	// ContainerMounts
	mounts := c.ContainerMounts("/cfg")
	if len(mounts) != 3 {
		t.Fatalf("ContainerMounts() returned %d mounts, want 3", len(mounts))
	}
	if mounts[0].Target != "/home/user/.claude" {
		t.Errorf("mounts[0].Target = %q, want /home/user/.claude", mounts[0].Target)
	}
	if mounts[1].Target != "/home/user/.claude.json" {
		t.Errorf("mounts[1].Target = %q, want /home/user/.claude.json", mounts[1].Target)
	}
	if mounts[2].Target != "/home/user/.config" {
		t.Errorf("mounts[2].Target = %q, want /home/user/.config", mounts[2].Target)
	}

	// GetDockerfileInstall
	df, err := c.GetDockerfileInstall("")
	if err != nil {
		t.Fatalf("GetDockerfileInstall() error: %v", err)
	}
	if !strings.Contains(df, "claude") {
		t.Error("GetDockerfileInstall() should reference claude")
	}

	// GetFullDockerfile
	full, err := c.GetFullDockerfile("1.0.0")
	if err != nil {
		t.Fatalf("GetFullDockerfile() error: %v", err)
	}
	if !strings.HasPrefix(full, "FROM exitbox-base") {
		t.Error("GetFullDockerfile() should start with FROM exitbox-base")
	}
}

func TestClaudeGetInstalledVersion_NilRuntime(t *testing.T) {
	c := &Claude{}
	if _, err := c.GetInstalledVersion(nil, "some-image"); err == nil {
		t.Errorf("GetInstalledVersion(nil, ...) should return error")
	}
}

func TestClaudeImportConfig(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	_ = os.WriteFile(filepath.Join(src, "settings.json"), []byte(`{"key":"val"}`), 0644)

	c := &Claude{}
	if err := c.ImportConfig(src, dst); err != nil {
		t.Fatalf("ImportConfig() error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dst, ".claude", "settings.json"))
	if err != nil {
		t.Fatalf("reading imported file: %v", err)
	}
	if string(data) != `{"key":"val"}` {
		t.Errorf("imported content = %q", string(data))
	}
}

func TestClaudeImportFile(t *testing.T) {
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
