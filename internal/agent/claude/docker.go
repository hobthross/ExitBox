package claude

func (c *Claude) GetFullDockerfile(version string) (string, error) {
	install, err := c.GetDockerfileInstall("")
	if err != nil {
		return "", err
	}
	return "FROM exitbox-base\n\n" + install, nil
}
func (c *Claude) GetDockerfileInstall(buildCtx string) (string, error) {
	return `# Install Claude Code via verified binary download (supply-chain hardened)
USER user
RUN set -e && \
    GCS_DEFAULT="` + claudeGCSDefault + `" && \
    INSTALL_SH_URL="` + claudeInstallSHURL + `" && \
    GCS_BUCKET="" && \
    if curl -fsSL --head "$GCS_DEFAULT/latest" >/dev/null 2>&1; then \
        GCS_BUCKET="$GCS_DEFAULT"; \
    else \
        echo "Hardcoded GCS URL unreachable, discovering from $INSTALL_SH_URL..." >&2; \
        DISCOVERED=$(curl -fsSL "$INSTALL_SH_URL" 2>/dev/null \
            | sed -n 's/^GCS_BUCKET="\(.*\)"/\1/p' | head -1) && \
        if [ -n "$DISCOVERED" ] && curl -fsSL --head "$DISCOVERED/latest" >/dev/null 2>&1; then \
            GCS_BUCKET="$DISCOVERED"; \
        fi; \
    fi && \
    if [ -z "$GCS_BUCKET" ]; then \
        echo "ERROR: Could not resolve Claude Code download URL" >&2; exit 1; \
    fi && \
    echo "Using GCS bucket: $GCS_BUCKET" && \
    case "$(uname -m)" in \
        x86_64|amd64) CLAUDE_ARCH="x64" ;; \
        aarch64|arm64) CLAUDE_ARCH="arm64" ;; \
        *) echo "Unsupported architecture: $(uname -m)" >&2; exit 1 ;; \
    esac && \
    CLAUDE_PLATFORM="linux-${CLAUDE_ARCH}-musl" && \
    CLAUDE_VERSION=$(curl -fsSL "$GCS_BUCKET/latest") && \
    echo "Installing Claude Code v${CLAUDE_VERSION} for ${CLAUDE_PLATFORM}..." && \
    MANIFEST=$(curl -fsSL "$GCS_BUCKET/$CLAUDE_VERSION/manifest.json") && \
    EXPECTED_CHECKSUM=$(printf '%s' "$MANIFEST" | jq -r ".platforms[\"$CLAUDE_PLATFORM\"].checksum // empty") && \
    if [ -z "$EXPECTED_CHECKSUM" ] || ! echo "$EXPECTED_CHECKSUM" | grep -qE '^[a-f0-9]{64}$'; then \
        echo "ERROR: No valid checksum for $CLAUDE_PLATFORM in manifest" >&2; exit 1; \
    fi && \
    mkdir -p "$HOME/.claude/downloads" && \
    BINARY_PATH="$HOME/.claude/downloads/claude-${CLAUDE_VERSION}" && \
    curl -fsSL -o "$BINARY_PATH" "$GCS_BUCKET/$CLAUDE_VERSION/$CLAUDE_PLATFORM/claude" && \
    ACTUAL_CHECKSUM=$(sha256sum "$BINARY_PATH" | cut -d' ' -f1) && \
    if [ "$ACTUAL_CHECKSUM" != "$EXPECTED_CHECKSUM" ]; then \
        echo "ERROR: Checksum verification failed!" >&2; \
        echo "  Expected: $EXPECTED_CHECKSUM" >&2; \
        echo "  Actual:   $ACTUAL_CHECKSUM" >&2; \
        rm -f "$BINARY_PATH"; exit 1; \
    fi && \
    echo "Checksum verified: $ACTUAL_CHECKSUM" && \
    chmod +x "$BINARY_PATH" && \
    "$BINARY_PATH" install && \
    rm -f "$BINARY_PATH" && \
    if [ -d "$HOME/.local/share/claude/versions" ]; then \
        latest_dir="$(ls -1d "$HOME/.local/share/claude/versions/"* | sort -V | tail -1)"; \
        if [ -x "$latest_dir/bin/claude" ]; then \
            ln -sf "$latest_dir/bin/claude" "$HOME/.local/bin/claude"; \
        fi; \
    fi && \
    command -v claude >/dev/null && \
    echo "Claude Code installed successfully"
USER root`, nil
}
