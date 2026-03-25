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
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/cloud-exit/exitbox/internal/container"
)

const qwenGitHubRepo = "QwenLM/qwen-code"

func (q *Qwen) GetLatestVersion() (string, error) {
	out, err := exec.Command("curl", "-s",
		fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", qwenGitHubRepo)).Output()
	if err != nil {
		return "", fmt.Errorf("failed to fetch Qwen Code latest version: %w", err)
	}
	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.Unmarshal(out, &release); err != nil {
		return "", err
	}
	v := strings.TrimPrefix(release.TagName, "v")
	if v == "" {
		return "", fmt.Errorf("empty tag_name")
	}
	return v, nil
}

func (q *Qwen) GetInstalledVersion(rt container.Runtime, img string) (string, error) {
	if rt == nil || !rt.ImageExists(img) {
		return "", fmt.Errorf("image %s not found", img)
	}
	out, err := rt.ImageInspect(img, `{{index .Config.Labels "exitbox.agent.version"}}`)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}
