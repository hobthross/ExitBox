package opencode

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestOpenCodeAgent(t *testing.T) {
	o := &OpenCode{}

	if o.Name() != "opencode" {
		t.Errorf("Name() = %q, want %q", o.Name(), "opencode")
	}
	if o.DisplayName() != "OpenCode" {
		t.Errorf("DisplayName() = %q, want %q", o.DisplayName(), "OpenCode")
	}

	// BinaryName
	bn := o.BinaryName()
	switch runtime.GOARCH {
	case "amd64":
		if bn != "opencode-linux-x64-musl.tar.gz" {
			t.Errorf("BinaryName() = %q on amd64", bn)
		}
	case "arm64":
		if bn != "opencode-linux-arm64-musl.tar.gz" {
			t.Errorf("BinaryName() = %q on arm64", bn)
		}
	}

	// HostConfigPaths
	paths := o.HostConfigPaths()
	if len(paths) != 2 {
		t.Fatalf("HostConfigPaths() returned %d paths, want 2", len(paths))
	}

	// ContainerMounts
	mounts := o.ContainerMounts("/cfg")
	if len(mounts) != 3 {
		t.Fatalf("ContainerMounts() returned %d mounts, want 3", len(mounts))
	}
	if mounts[0].Target != "/home/user/.opencode" {
		t.Errorf("mounts[0].Target = %q, want /home/user/.opencode", mounts[0].Target)
	}

	// GetDockerfileInstall
	df, err := o.GetDockerfileInstall("")
	if err != nil {
		t.Fatalf("GetDockerfileInstall() error: %v", err)
	}
	if !strings.Contains(df, "sha256sum") {
		t.Error("GetDockerfileInstall() should contain sha256sum verification")
	}

	// GetFullDockerfile
	full, err := o.GetFullDockerfile("0.2.0")
	if err != nil {
		t.Fatalf("GetFullDockerfile() error: %v", err)
	}
	if !strings.HasPrefix(full, "FROM exitbox-base") {
		t.Error("GetFullDockerfile() should start with FROM exitbox-base")
	}
	if !strings.Contains(full, "OPENCODE_VERSION=0.2.0") {
		t.Error("GetFullDockerfile() should include OPENCODE_VERSION ARG")
	}
}

func TestOpenCodeGetInstalledVersion_NilRuntime(t *testing.T) {
	o := &OpenCode{}
	if _, err := o.GetInstalledVersion(nil, "some-image"); err == nil {
		t.Errorf("GetInstalledVersion(nil, ...) should return error")
	}
}

func TestOpenCodeImportConfig_DefaultDir(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	_ = os.WriteFile(filepath.Join(src, "config.json"), []byte(`{}`), 0644)

	o := &OpenCode{}
	if err := o.ImportConfig(src, dst); err != nil {
		t.Fatalf("ImportConfig() error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dst, ".opencode", "config.json")); err != nil {
		t.Errorf("expected .opencode/config.json to exist: %v", err)
	}
}

func TestOpenCodeImportConfig_ConfigDir(t *testing.T) {
	src := filepath.Join(t.TempDir(), ".config", "opencode")
	dst := t.TempDir()
	_ = os.MkdirAll(src, 0755)
	_ = os.WriteFile(filepath.Join(src, "settings.json"), []byte(`{}`), 0644)

	o := &OpenCode{}
	if err := o.ImportConfig(src, dst); err != nil {
		t.Fatalf("ImportConfig() error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dst, ".config", "opencode", "settings.json")); err != nil {
		t.Errorf("expected .config/opencode/settings.json to exist: %v", err)
	}
}


