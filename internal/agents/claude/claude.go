package claude

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloud-exit/exitbox/internal/agent"
	"github.com/cloud-exit/exitbox/internal/config"
	"github.com/cloud-exit/exitbox/internal/fsutil"
)

const (
	claudeGCSDefault   = "https://storage.googleapis.com/claude-code-dist-86c565f3-f756-42ad-8dfa-d59b1c096819/claude-code-releases"
	claudeInstallSHURL = "https://claude.ai/install.sh"
)

type Claude struct{}

var _ agent.Agent = (*Claude)(nil)

func (c *Claude) Name() string        { return "claude" }
func (c *Claude) DisplayName() string { return "Claude Code" }
func (c *Claude) Description() string { return "Anthropic's AI coding assistant" }

func (c *Claude) OllamaEnvVars(ollamaBaseURL string) []string {
	return []string{
		"ANTHROPIC_BASE_URL=" + ollamaBaseURL,
		"ANTHROPIC_AUTH_TOKEN=ollama",
		"ANTHROPIC_API_KEY=",
	}
}

func (c *Claude) HostConfigPaths() []string {
	home := os.Getenv("HOME")
	return []string{filepath.Join(home, ".claude")}
}

func (c *Claude) ContainerMounts(cfgDir string) []agent.Mount {
	return []agent.Mount{
		{Source: filepath.Join(cfgDir, ".claude"), Target: "/home/user/.claude"},
		{Source: filepath.Join(cfgDir, ".claude.json"), Target: "/home/user/.claude.json"},
		{Source: filepath.Join(cfgDir, ".config"), Target: "/home/user/.config"},
	}
}

func (c *Claude) DetectHostConfig() (string, error) {
	home := os.Getenv("HOME")
	dir := filepath.Join(home, ".claude")
	if info, err := os.Stat(dir); err == nil && info.IsDir() {
		return dir, nil
	}
	return "", fmt.Errorf("no Claude config found")
}

func (c *Claude) EnsureWorkspaceAgentConfig(workspaceName string) error {
	if workspaceName == "" {
		return nil
	}
	root := config.WorkspaceAgentDir(workspaceName, c.Name())
	_ = os.MkdirAll(root, 0755)
	home := os.Getenv("HOME")

	claudeDir := fsutil.EnsureDir(root, ".claude")
	fsutil.SeedDirOnce(filepath.Join(home, ".claude"), claudeDir)

	claudeJSON := fsutil.EnsureFile(root, ".claude.json")
	fsutil.SeedFileOnce(filepath.Join(home, ".claude.json"), claudeJSON)

	cfgDir := fsutil.EnsureDir(root, ".config")
	fsutil.SeedDirOnce(filepath.Join(home, ".config"), cfgDir)
	return nil
}

func (c *Claude) ImportConfig(src, dst string) error {
	home := os.Getenv("HOME")

	// Copy entire ~/.claude directory
	target := filepath.Join(dst, ".claude")
	_ = os.MkdirAll(target, 0755)
	if err := agent.CopyDir(src, target); err != nil {
		return fmt.Errorf("copying .claude dir: %w", err)
	}

	// Also copy ~/.claude.json if it exists
	claudeJSON := filepath.Join(home, ".claude.json")
	if data, err := os.ReadFile(claudeJSON); err == nil {
		_ = os.WriteFile(filepath.Join(dst, ".claude.json"), data, 0644)
	}

	return nil
}
