package codex

import (
	"fmt"
	"strings"
)

func (c *Codex) GetDockerfileInstall(buildCtx string) (string, error) {
	binaryName := c.BinaryName()
	if binaryName == "" {
		return "", fmt.Errorf("unsupported architecture for Codex")
	}
	binaryInside := strings.TrimSuffix(binaryName, ".tar.gz")

	return fmt.Sprintf(`# Install Codex binary with SHA-256 verification
ARG CODEX_VERSION
ARG CODEX_CHECKSUM
COPY %s /tmp/codex.tar.gz
RUN echo "${CODEX_CHECKSUM}  /tmp/codex.tar.gz" | sha256sum -c - && \
    mkdir -p $HOME/.local/bin && \
    tar -xzf /tmp/codex.tar.gz -C /tmp && \
    mv /tmp/%s $HOME/.local/bin/codex && \
    chmod +x $HOME/.local/bin/codex && \
    rm -f /tmp/codex.tar.gz && \
    $HOME/.local/bin/codex --version`, binaryName, binaryInside), nil
}

func (c *Codex) GetFullDockerfile(version string) (string, error) {
	install, err := c.GetDockerfileInstall("")
	if err != nil {
		return "", err
	}
	df := "FROM exitbox-base\n\n"
	if version != "" {
		df += fmt.Sprintf("ARG CODEX_VERSION=%s\n", version)
	}
	df += install
	return df, nil
}
