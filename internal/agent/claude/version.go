package claude

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/cloud-exit/exitbox/internal/container"
)

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