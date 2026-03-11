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
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cloud-exit/exitbox/internal/container"
)

const (
	claudeGCSDefault   = "https://storage.googleapis.com/claude-code-dist-86c565f3-f756-42ad-8dfa-d59b1c096819/claude-code-releases"
	claudeInstallSHURL = "https://claude.ai/install.sh"
)

// Claude implements the Agent interface for Claude Code.
type Claude struct{}

func (c *Claude) Name() string        { return "claude" }
func (c *Claude) DisplayName() string { return "Claude Code" }

func (c *Claude) GetLatestVersion() (string, error) {
	out, err := exec.Command("curl", "-fsSL", claudeGCSDefault+"/latest").Output()
	if err != nil {
		return "", fmt.Errorf("failed to fetch Claude latest version: %w", err)
	}
	v := strings.TrimSpace(string(out))
	if v == "" {
		return "", fmt.Errorf("empty version response")
	}
	return v, nil
}

func (c *Claude) GetInstalledVersion(rt container.Runtime, img string) (string, error) {
	if rt == nil || !rt.ImageExists(img) {
		return "", fmt.Errorf("image %s not found", img)
	}
	out, err := rt.ImageInspect(img, `{{index .Config.Labels "exitbox.agent.version"}}`)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func (c *Claude) GetDockerfileInstall(buildCtx string) (string, error) {
	return `
RUN apk add --no-cache musl-dev gcc && \
    printf '#include <sys/syscall.h>\n#include <unistd.h>\nint posix_getdents(int fd, void *buf, unsigned long nbytes, int flags) {\n  (void)flags;\n  return syscall(SYS_getdents64, fd, buf, nbytes);\n}\n' \
        > /tmp/posix_getdents.c && \
    gcc -shared -o /usr/local/lib/posix_getdents.so /tmp/posix_getdents.c && \
    rm /tmp/posix_getdents.c && \
    apk del musl-dev gcc

ENV LD_PRELOAD="/lib/libgcompat.so.0 /usr/local/lib/posix_getdents.so"

ARG CLAUDE_VERSION
RUN set -e && \
    if [ -n "$CLAUDE_VERSION" ]; then \
        curl -fsSL https://claude.ai/install.sh | bash -s "$CLAUDE_VERSION"; \
    else \
        curl -fsSL https://claude.ai/install.sh | bash; \
    fi && \
    claude --version && \
    echo "Claude Code installed successfully"

# Configure settings for runtime user (as root, since /home/user may not be owned by user yet)
RUN mkdir -p /home/user/.claude && \
    echo '{"env":{"DISABLE_AUTOUPDATER":"1","USE_BUILTIN_RIPGREP":"0"}}' > /home/user/.claude/settings.json && \
    chown -R user:user /home/user/.claude`, nil
}



func (c *Claude) GetFullDockerfile(version string) (string, error) {
	install, err := c.GetDockerfileInstall("")
	if err != nil {
		return "", err
	}
	// Add CLAUDE_VERSION build arg if version is specified
	if version != "" {
		return fmt.Sprintf("FROM exitbox-base\n\nARG CLAUDE_VERSION=%s\n\n", version) + install, nil
	}
	return "FROM exitbox-base\n\n" + install, nil
}

func (c *Claude) HostConfigPaths() []string {
	home := os.Getenv("HOME")
	return []string{filepath.Join(home, ".claude")}
}

func (c *Claude) ContainerMounts(cfgDir string) []Mount {
	return []Mount{
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

func (c *Claude) ImportFile(src, dst string) error {
	target := filepath.Join(dst, ".claude", filepath.Base(src))
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return err
	}
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(target, data, 0644)
}

func (c *Claude) ConfigFilePath(wsDir string) string {
	return filepath.Join(wsDir, ".claude", "settings.json")
}

func (c *Claude) ImportConfig(src, dst string) error {
	home := os.Getenv("HOME")

	// Copy entire ~/.claude directory
	target := filepath.Join(dst, ".claude")
	_ = os.MkdirAll(target, 0755)
	if err := copyDirContents(src, target); err != nil {
		return fmt.Errorf("copying .claude dir: %w", err)
	}

	// Also copy ~/.claude.json if it exists
	claudeJSON := filepath.Join(home, ".claude.json")
	if data, err := os.ReadFile(claudeJSON); err == nil {
		_ = os.WriteFile(filepath.Join(dst, ".claude.json"), data, 0644)
	}

	return nil
}
