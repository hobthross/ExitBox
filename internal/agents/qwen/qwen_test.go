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

package qwen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestQwenAgent(t *testing.T) {
	q := &Qwen{}

	if q.Name() != "qwen" {
		t.Errorf("Name() = %q, want %q", q.Name(), "qwen")
	}
	if q.DisplayName() != "Qwen Code" {
		t.Errorf("DisplayName() = %q, want %q", q.DisplayName(), "Qwen Code")
	}

	paths := q.HostConfigPaths()
	if len(paths) != 2 {
		t.Fatalf("HostConfigPaths() returned %d paths, want 2", len(paths))
	}
	if !strings.HasSuffix(paths[0], ".qwen") {
		t.Errorf("HostConfigPaths()[0] = %q, want path ending in .qwen", paths[0])
	}
	if !strings.Contains(paths[1], filepath.Join(".config", "qwen")) {
		t.Errorf("HostConfigPaths()[1] = %q, want path containing .config/qwen", paths[1])
	}

	mounts := q.ContainerMounts("/cfg")
	if len(mounts) != 2 {
		t.Fatalf("ContainerMounts() returned %d mounts, want 2", len(mounts))
	}
	if mounts[0].Target != "/home/user/.qwen" {
		t.Errorf("mounts[0].Target = %q, want /home/user/.qwen", mounts[0].Target)
	}
	if mounts[1].Target != "/home/user/.config/qwen" {
		t.Errorf("mounts[1].Target = %q, want /home/user/.config/qwen", mounts[1].Target)
	}

	df, err := q.GetDockerfileInstall("")
	if err != nil {
		t.Fatalf("GetDockerfileInstall() error: %v", err)
	}
	if !strings.Contains(df, "npm install") {
		t.Error("GetDockerfileInstall() should install via npm")
	}
	if !strings.Contains(df, "@qwen-code/qwen-code") {
		t.Error("GetDockerfileInstall() should reference @qwen-code/qwen-code")
	}

	full, err := q.GetFullDockerfile("0.11.0")
	if err != nil {
		t.Fatalf("GetFullDockerfile() error: %v", err)
	}
	if !strings.HasPrefix(full, "FROM exitbox-base") {
		t.Error("GetFullDockerfile() should start with FROM exitbox-base")
	}
	if !strings.Contains(full, "QWEN_VERSION=0.11.0") {
		t.Error("GetFullDockerfile() should include QWEN_VERSION ARG")
	}
}

func TestQwenGetInstalledVersion_NilRuntime(t *testing.T) {
	q := &Qwen{}
	if _, err := q.GetInstalledVersion(nil, "some-image"); err == nil {
		t.Errorf("GetInstalledVersion(nil, ...) should return error")
	}
}

func TestQwenImportConfig_DefaultDir(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	_ = os.MkdirAll(filepath.Join(src, "nested"), 0755)
	_ = os.WriteFile(filepath.Join(src, "settings.json"), []byte(`{}`), 0644)

	q := &Qwen{}
	if err := q.ImportConfig(src, dst); err != nil {
		t.Fatalf("ImportConfig() error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dst, ".qwen", "settings.json")); err != nil {
		t.Errorf("expected .qwen/settings.json to exist: %v", err)
	}
}

func TestQwenImportConfig_ConfigDir(t *testing.T) {
	src := filepath.Join(t.TempDir(), ".config", "qwen")
	dst := t.TempDir()
	_ = os.MkdirAll(src, 0755)
	_ = os.WriteFile(filepath.Join(src, "settings.json"), []byte(`{}`), 0644)

	q := &Qwen{}
	if err := q.ImportConfig(src, dst); err != nil {
		t.Fatalf("ImportConfig() error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dst, ".config", "qwen", "settings.json")); err != nil {
		t.Errorf("expected .config/qwen/settings.json to exist: %v", err)
	}
}
