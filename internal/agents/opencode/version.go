package opencode

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/cloud-exit/exitbox/internal/container"
)

const opencodeGitHubRepo = "openai/opencode"

func (o *OpenCode) GetLatestVersion() (string, error) {
	out, err := exec.Command("curl", "-s",
		fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", opencodeGitHubRepo)).Output()
	if err != nil {
		return "", fmt.Errorf("failed to fetch OpenCode latest version: %w", err)
	}
	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.Unmarshal(out, &release); err != nil {
		return "", err
	}
	// Strip leading 'v' if present
	v := strings.TrimPrefix(release.TagName, "v")
	if v == "" {
		return "", fmt.Errorf("empty tag_name")
	}
	return v, nil
}

func (o *OpenCode) GetInstalledVersion(rt container.Runtime, img string) (string, error) {
	if rt == nil || !rt.ImageExists(img) {
		return "", fmt.Errorf("image %s not found", img)
	}
	out, err := rt.ImageInspect(img, `{{index .Config.Labels "exitbox.agent.version"}}`)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}
