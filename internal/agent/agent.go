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
	"context"

	"github.com/cloud-exit/exitbox/internal/config"
	"github.com/cloud-exit/exitbox/internal/container"
)

// Mount describes a volume mount for a container.
type Mount struct {
	Source string
	Target string
}

// PrepareBuildInput holds parameters for preparing an agent's Docker build context.
// It is passed to Agent.PrepareBuild to support download, checksum, and logging.
type PrepareBuildInput struct {
	Ctx            context.Context
	Version        string
	BuildDir       string
	DockerfilePath string
	Download       func(ctx context.Context, url, dest string) error
	FileSHA256     func(path string) string
	Logf           func(format string, args ...interface{})
}

type Agent interface {
	Name() string
	DisplayName() string
	Description() string

	GetLatestVersion() (string, error)
	GetInstalledVersion(rt container.Runtime, img string) (string, error)

	GetDockerfileInstall(buildCtx string) (string, error)
	GetFullDockerfile(version string) (string, error)

	HostConfigPaths() []string
	ContainerMounts(cfgDir string) []Mount
	DetectHostConfig() (string, error)
	ImportConfig(src, dst string) error

	GenerateConfig(cfg config.ServerConfig) (map[string]interface{}, error)
	LogSearchDirs(home, agentCfgDir string) []string

	// PrepareBuild prepares the build context (downloads binaries, writes Dockerfile) for the agent's core image.
	PrepareBuild(in PrepareBuildInput) error

	// EnsureWorkspaceAgentConfig creates and seeds agent config dirs for the given workspace.
	EnsureWorkspaceAgentConfig(workspaceName string) error

	OllamaEnvVars(ollamaBaseURL string) []string
	ConfigFilePath(agentDir string) string
	// ExtractConfigServerURLs returns server URLs from parsed agent config JSON.
	ExtractConfigServerURLs(data map[string]interface{}) []string
}
