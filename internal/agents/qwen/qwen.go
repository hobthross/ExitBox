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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloud-exit/exitbox/internal/agent"
	"github.com/cloud-exit/exitbox/internal/config"
	"github.com/cloud-exit/exitbox/internal/fsutil"
)

type Qwen struct{}

var _ agent.Agent = (*Qwen)(nil)

func (q *Qwen) Name() string        { return "qwen" }
func (q *Qwen) DisplayName() string { return "Qwen Code" }
func (q *Qwen) Description() string {
	return "Open-source AI agent for the terminal (Qwen Code)"
}

func (q *Qwen) OllamaEnvVars(ollamaBaseURL string) []string {
	return []string{
		"OPENAI_BASE_URL=" + ollamaBaseURL + "/v1",
		"OPENAI_API_KEY=ollama",
	}
}

func (q *Qwen) HostConfigPaths() []string {
	home := os.Getenv("HOME")
	return []string{
		filepath.Join(home, ".qwen"),
		filepath.Join(home, ".config", "qwen"),
	}
}

func (q *Qwen) ContainerMounts(cfgDir string) []agent.Mount {
	return []agent.Mount{
		{Source: filepath.Join(cfgDir, ".qwen"), Target: "/home/user/.qwen"},
		{Source: filepath.Join(cfgDir, ".config", "qwen"), Target: "/home/user/.config/qwen"},
	}
}

func (q *Qwen) EnsureWorkspaceAgentConfig(workspaceName string) error {
	if workspaceName == "" {
		return nil
	}
	root := config.WorkspaceAgentDir(workspaceName, q.Name())
	_ = os.MkdirAll(root, 0755)
	home := os.Getenv("HOME")

	qwenDir := fsutil.EnsureDir(root, ".qwen")
	fsutil.SeedDirOnce(filepath.Join(home, ".qwen"), qwenDir)

	qwenCfg := fsutil.EnsureDir(root, ".config", "qwen")
	fsutil.SeedDirOnce(filepath.Join(home, ".config", "qwen"), qwenCfg)

	return nil
}

func (q *Qwen) DetectHostConfig() (string, error) {
	home := os.Getenv("HOME")
	for _, p := range []string{
		filepath.Join(home, ".qwen"),
		filepath.Join(home, ".config", "qwen"),
	} {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return p, nil
		}
	}
	return "", fmt.Errorf("no Qwen config found")
}

func (q *Qwen) ImportConfig(src, dst string) error {
	if strings.Contains(src, filepath.Join(".config", "qwen")) {
		target := filepath.Join(dst, ".config", "qwen")
		_ = os.MkdirAll(target, 0755)
		return fsutil.CopyDir(src, target)
	}
	target := filepath.Join(dst, ".qwen")
	_ = os.MkdirAll(target, 0755)
	return fsutil.CopyDir(src, target)
}

func (q *Qwen) ImportFile(src, dst string) error {
	target := filepath.Join(dst, ".qwen", filepath.Base(src))
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return err
	}
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(target, data, 0644)
}
