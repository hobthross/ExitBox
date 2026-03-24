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

package opencode

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloud-exit/exitbox/internal/agent"
)

// opencodeReleaseRepo is the GitHub org/repo for OpenCode release downloads (v-prefixed tags).
const opencodeReleaseRepo = "anomalyco/opencode"

func (o *OpenCode) GetDockerfileInstall(buildCtx string) (string, error) {
	return fmt.Sprintf(`# Install OpenCode binary with SHA-256 verification
ARG OPENCODE_VERSION
ARG OPENCODE_CHECKSUM
COPY %s /tmp/opencode.tar.gz
RUN echo "${OPENCODE_CHECKSUM}  /tmp/opencode.tar.gz" | sha256sum -c - && \
    tar -xzf /tmp/opencode.tar.gz -C /usr/local/bin && \
    chmod +x /usr/local/bin/opencode && \
    rm -f /tmp/opencode.tar.gz`, o.BinaryName()), nil
}

// GetFullDockerfile returns the complete Dockerfile for OpenCode.
// Builds on exitbox-base with pre-downloaded musl binary (same pattern as Claude/Codex).
func (o *OpenCode) GetFullDockerfile(version string) (string, error) {
	install, err := o.GetDockerfileInstall("")
	if err != nil {
		return "", err
	}
	df := "FROM exitbox-base\n\n"
	if version != "" {
		df += fmt.Sprintf("ARG OPENCODE_VERSION=%s\n", version)
	}
	df += install
	return df, nil
}

func (o *OpenCode) PrepareBuild(in agent.PrepareBuildInput) error {
	version := in.Version
	if version == "" {
		version = "latest"
	}
	binaryName := o.BinaryName()
	if binaryName == "" {
		return fmt.Errorf("unsupported architecture for OpenCode")
	}
	if in.Download == nil || in.FileSHA256 == nil {
		return fmt.Errorf("PrepareBuildInput.Download and FileSHA256 are required for OpenCode")
	}
	url := fmt.Sprintf("https://github.com/%s/releases/download/v%s/%s", opencodeReleaseRepo, version, binaryName)
	if in.Logf != nil {
		in.Logf("Downloading OpenCode %s...", version)
	}
	dlPath := filepath.Join(in.BuildDir, binaryName)
	if err := in.Download(in.Ctx, url, dlPath); err != nil {
		return fmt.Errorf("failed to download OpenCode: %w", err)
	}
	checksum := in.FileSHA256(dlPath)
	if in.Logf != nil {
		in.Logf("OpenCode SHA-256: %s", checksum)
	}
	df := fmt.Sprintf("FROM exitbox-base\n\nARG OPENCODE_VERSION=%s\nARG OPENCODE_CHECKSUM=%s\n", version, checksum)
	install, err := o.GetDockerfileInstall(in.BuildDir)
	if err != nil {
		return fmt.Errorf("failed to get OpenCode install instructions: %w", err)
	}
	df += install
	if err := os.WriteFile(in.DockerfilePath, []byte(df), 0644); err != nil {
		return fmt.Errorf("failed to write Dockerfile: %w", err)
	}
	return nil
}
