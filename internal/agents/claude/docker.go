package claude

import (
	"fmt"
	"os"

	"github.com/cloud-exit/exitbox/internal/agent"
)

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
func (c *Claude) GetDockerfileInstall(buildCtx string) (string, error) {
	return `# Install Claude Code via verified binary download (supply-chain hardened)
ARG CLAUDE_VERSION
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
    if [ -n "$CLAUDE_VERSION" ]; then \
        TRY_VERSION="$CLAUDE_VERSION"; \
    else \
        TRY_VERSION=$(curl -fsSL "$GCS_BUCKET/latest"); \
    fi && \
    claude_try_install() { \
        local ver="$1"; \
        echo "Installing Claude Code v${ver} for ${CLAUDE_PLATFORM}..." && \
        MANIFEST=$(curl -fsSL "$GCS_BUCKET/$ver/manifest.json") && \
        EXPECTED_CHECKSUM=$(printf '%s' "$MANIFEST" | jq -r ".platforms[\"$CLAUDE_PLATFORM\"].checksum // empty") && \
        if [ -z "$EXPECTED_CHECKSUM" ] || ! echo "$EXPECTED_CHECKSUM" | grep -qE '^[a-f0-9]{64}$'; then \
            echo "No valid checksum for $CLAUDE_PLATFORM in v${ver} manifest" >&2; return 1; \
        fi && \
        mkdir -p "$HOME/.claude/downloads" && \
        local bp="$HOME/.claude/downloads/claude-${ver}" && \
        curl -fsSL -o "$bp" "$GCS_BUCKET/$ver/$CLAUDE_PLATFORM/claude" && \
        ACTUAL_CHECKSUM=$(sha256sum "$bp" | cut -d' ' -f1) && \
        if [ "$ACTUAL_CHECKSUM" != "$EXPECTED_CHECKSUM" ]; then \
            echo "Checksum mismatch for v${ver}" >&2; rm -f "$bp"; return 1; \
        fi && \
        echo "Checksum verified: $ACTUAL_CHECKSUM" && \
        chmod +x "$bp" && \
        if ! "$bp" --version >/dev/null 2>&1; then \
            echo "Binary v${ver} failed compatibility check (musl)" >&2; rm -f "$bp"; return 1; \
        fi && \
        # Skip the built-in install command which may trigger auto-update
        # Instead, manually copy the binary to the target location
        mkdir -p "$HOME/.local/share/claude/versions/$ver/bin" && \
        cp "$bp" "$HOME/.local/share/claude/versions/$ver/bin/claude" && \
        chmod +x "$HOME/.local/share/claude/versions/$ver/bin/claude" && \
        rm -f "$bp"; \
    } && \
    INSTALLED=false && \
    if claude_try_install "$TRY_VERSION"; then \
        INSTALLED=true; \
    elif [ -z "$CLAUDE_VERSION" ]; then \
        MAJOR=$(echo "$TRY_VERSION" | cut -d. -f1); \
        MINOR=$(echo "$TRY_VERSION" | cut -d. -f2); \
        PATCH=$(echo "$TRY_VERSION" | cut -d. -f3); \
        FALLBACK=0; \
        while [ "$INSTALLED" = "false" ] && [ "$FALLBACK" -lt 5 ] && [ "$PATCH" -gt 0 ]; do \
            PATCH=$((PATCH - 1)); \
            FALLBACK=$((FALLBACK + 1)); \
            PREV="${MAJOR}.${MINOR}.${PATCH}"; \
            echo "Trying fallback v${PREV}..." >&2; \
            if claude_try_install "$PREV"; then \
                INSTALLED=true; \
            fi; \
        done; \
    fi && \
    if [ "$INSTALLED" = "false" ]; then \
        echo "ERROR: Could not install any compatible Claude Code version" >&2; exit 1; \
    fi && \
		# Only symlink if we didn't already install a specific version
		if [ -n "$CLAUDE_VERSION" ]; then \
			# When version is pinned, symlink the specific version we just installed
			INSTALLED_DIR="$HOME/.local/share/claude/versions/$TRY_VERSION"; \
			if [ -d "$INSTALLED_DIR" ] && [ -x "$INSTALLED_DIR/bin/claude" ]; then \
				ln -sf "$INSTALLED_DIR/bin/claude" "$HOME/.local/bin/claude"; \
			fi; \
		elif [ -d "$HOME/.local/share/claude/versions" ]; then \
			latest_dir="$(ls -1d "$HOME/.local/share/claude/versions/"* | sort -V | tail -1)"; \
			if [ -x "$latest_dir/bin/claude" ]; then \
				ln -sf "$latest_dir/bin/claude" "$HOME/.local/bin/claude"; \
			fi; \
		fi && \
		# Disable auto-updater by setting the environment variable AND writing to settings.json
		export DISABLE_AUTOUPDATER=1 && \
		mkdir -p "$HOME/.claude" && \
		if [ -f "$HOME/.claude/settings.json" ]; then \
			# Add DISABLE_AUTOUPDATER to existing settings.json
			if ! grep -q "DISABLE_AUTOUPDATER" "$HOME/.claude/settings.json"; then \
				sed -i 's/}$/,"env":{"DISABLE_AUTOUPDATER":"1"}\n}/' "$HOME/.claude/settings.json" 2>/dev/null || \
				echo '{"env":{"DISABLE_AUTOUPDATER":"1"}}' >> "$HOME/.claude/settings.json"; \
			fi; \
		else \
			# Create settings.json with auto-updater disabled
			echo '{"env":{"DISABLE_AUTOUPDATER":"1"}}' > "$HOME/.claude/settings.json"; \
		fi && \
		command -v claude >/dev/null && \
		echo "Claude Code installed successfully (auto-updater disabled)"
USER root`, nil
}

func (c *Claude) PrepareBuild(in agent.PrepareBuildInput) error {
	df, err := c.GetFullDockerfile(in.Version)
	if err != nil {
		return err
	}
	if err := os.WriteFile(in.DockerfilePath, []byte(df), 0644); err != nil {
		return fmt.Errorf("failed to write Dockerfile: %w", err)
	}
	return nil
}
