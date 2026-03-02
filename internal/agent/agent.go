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

	"github.com/cloud-exit/exitbox/internal/container"
	"github.com/cloud-exit/exitbox/internal/generate"
)

// Mount describes a volume mount for a container.
type Mount struct {
	Source string
	Target string
}

type ConfigGenerator interface {
	GenerateConfig(cfg generate.ServerConfig) (map[string]interface{}, error)
}

type LogLocationProvider interface {
	LogSearchDirs(home, agentCfgDir string) []string
}

type PrepareBuildProvider interface {
	// PrepareBuild prepares the build context (downloads binaries, writes Dockerfile) for the agent's core image.
	PrepareBuild(in PrepareBuildInput) error
}

// WorkspaceConfigEnsurer ensures workspace-specific agent config directories exist
// and seeds them from host config when empty.
type WorkspaceConfigEnsurer interface {
	// EnsureWorkspaceAgentConfig creates and seeds agent config dirs for the given workspace.
	EnsureWorkspaceAgentConfig(workspaceName string) error
}

// PrepareBuildInput holds parameters for preparing an agent's Docker build context.
// It is passed to Agent.PrepareBuild to support download, checksum, and logging.
type PrepareBuildInput struct {
	// Ctx is the context for cancellable operations (e.g. HTTP downloads).
	Ctx context.Context
	// Version is the agent version to build (e.g. "v1.0.0" or "latest").
	Version string
	// BuildDir is the build context directory (artifacts like pre-downloaded binaries go here).
	BuildDir string
	// DockerfilePath is the path where the Dockerfile must be written.
	DockerfilePath string
	// Download performs an HTTP GET and writes the response to dest. May be nil if agent needs no downloads.
	Download func(ctx context.Context, url, dest string) error
	// FileSHA256 returns the SHA-256 hex digest of the file at path. May be nil if agent needs no checksums.
	FileSHA256 func(path string) string
	// Logf logs an informational message. May be nil (logging is skipped).
	Logf func(format string, args ...interface{})
}

// Agent is the interface that all agent implementations must satisfy.
type Agent interface {
	Name() string
	DisplayName() string
	GetLatestVersion() (string, error)
	GetInstalledVersion(rt container.Runtime, img string) (string, error)
	GetDockerfileInstall(buildCtx string) (string, error)
	GetFullDockerfile(version string) (string, error)
	HostConfigPaths() []string
	ContainerMounts(cfgDir string) []Mount
	DetectHostConfig() (string, error)
	ImportConfig(src, dst string) error
}

type AgentEntity interface {
	Agent
	ConfigGenerator
	LogLocationProvider
	PrepareBuildProvider
	WorkspaceConfigEnsurer
}

// AgentNames is the list of all supported agent names.
var AgentNames = []string{"claude", "codex", "opencode"}

// DisplayName returns the human-readable name for an agent.
func DisplayName(name string) string {
	switch name {
	case "claude":
		return "Claude Code"
	case "codex":
		return "OpenAI Codex"
	case "opencode":
		return "OpenCode"
	}
	return name
}

// IsValidAgent returns true if the name is a known agent.
func IsValidAgent(name string) bool {
	for _, a := range AgentNames {
		if a == name {
			return true
		}
	}
	return false
}

// registry holds all agent implementations.
var registry = map[string]Agent{}

// Register adds an agent to the registry.
func Register(a Agent) {
	registry[a.Name()] = a
}

// Get returns the agent implementation for a name.
func Get(name string) Agent {
	return registry[name]
}

func init() {
	Register(&Claude{})
	Register(&Codex{})
	Register(&OpenCode{})
}
