package opencode

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cloud-exit/exitbox/internal/agent"
	"github.com/cloud-exit/exitbox/internal/config"
	"github.com/cloud-exit/exitbox/internal/fsutil"
)

type OpenCode struct{}

var _ agent.AgentEntity = (*OpenCode)(nil)

func (o *OpenCode) Name() string        { return "opencode" }
func (o *OpenCode) DisplayName() string { return "OpenCode" }

// BinaryName returns the platform-specific binary tarball name (musl build for Alpine).
func (o *OpenCode) BinaryName() string {
	arch := runtime.GOARCH
	switch arch {
	case "amd64":
		return "opencode-linux-x64-musl.tar.gz"
	case "arm64":
		return "opencode-linux-arm64-musl.tar.gz"
	default:
		return ""
	}
}

func (o *OpenCode) HostConfigPaths() []string {
	home := os.Getenv("HOME")
	return []string{
		filepath.Join(home, ".opencode"),
		filepath.Join(home, ".config", "opencode"),
	}
}

func (o *OpenCode) ContainerMounts(cfgDir string) []agent.Mount {
	return []agent.Mount{
		{Source: filepath.Join(cfgDir, ".opencode"), Target: "/home/user/.opencode"},
		{Source: filepath.Join(cfgDir, ".config", "opencode"), Target: "/home/user/.config/opencode"},
		{Source: filepath.Join(cfgDir, ".local", "share", "opencode"), Target: "/home/user/.local/share/opencode"},
	}
}

func (o *OpenCode) EnsureWorkspaceAgentConfig(workspaceName string) error {
	if workspaceName == "" {
		return nil
	}
	root := config.WorkspaceAgentDir(workspaceName, o.Name())
	_ = os.MkdirAll(root, 0755)
	home := os.Getenv("HOME")

	ocDir := fsutil.EnsureDir(root, ".opencode")
	fsutil.SeedDirOnce(filepath.Join(home, ".opencode"), ocDir)

	ocCfg := fsutil.EnsureDir(root, ".config", "opencode")
	fsutil.SeedDirOnce(filepath.Join(home, ".config", "opencode"), ocCfg)

	ocShare := fsutil.EnsureDir(root, ".local", "share", "opencode")
	fsutil.SeedDirOnce(filepath.Join(home, ".local", "share", "opencode"), ocShare)

	ocState := fsutil.EnsureDir(root, ".local", "state")
	fsutil.SeedDirOnce(filepath.Join(home, ".local", "state"), ocState)

	ocCache := fsutil.EnsureDir(root, ".cache", "opencode")
	fsutil.SeedDirOnce(filepath.Join(home, ".cache", "opencode"), ocCache)
	return nil
}

func (o *OpenCode) DetectHostConfig() (string, error) {
	home := os.Getenv("HOME")
	for _, p := range []string{
		filepath.Join(home, ".opencode"),
		filepath.Join(home, ".config", "opencode"),
	} {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return p, nil
		}
	}
	return "", fmt.Errorf("no OpenCode config found")
}

func (o *OpenCode) ImportConfig(src, dst string) error {
	if strings.Contains(src, filepath.Join(".config", "opencode")) {
		target := filepath.Join(dst, ".config", "opencode")
		_ = os.MkdirAll(target, 0755)
		return agent.CopyDir(src, target)
	}
	target := filepath.Join(dst, ".opencode")
	_ = os.MkdirAll(target, 0755)
	return agent.CopyDir(src, target)
}
