package opencode

import "fmt"

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
