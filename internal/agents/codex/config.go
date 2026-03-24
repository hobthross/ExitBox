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

package codex

import (
	"path/filepath"

	"github.com/cloud-exit/exitbox/internal/config"
)

// GenerateConfig produces a Codex config.json config map.
func (c *Codex) GenerateConfig(cfg config.ServerConfig) (map[string]interface{}, error) {
	return map[string]interface{}{
		"model":    cfg.ProviderID + "/" + cfg.ModelID,
		"provider": cfg.BaseURL,
	}, nil
}

func (c *Codex) LogSearchDirs(home, agentCfgDir string) []string {
	return []string{
		filepath.Join(home, ".codex"),
		filepath.Join(home, ".config", "codex"),
		filepath.Join(agentCfgDir, ".codex"),
		filepath.Join(agentCfgDir, ".config", "codex"),
	}
}

func (c *Codex) ConfigFilePath(agentDir string) string {
	return filepath.Join(agentDir, ".codex", "config.json")
}

// ExtractConfigServerURLs reads the top-level provider URL from a Codex config.
func (c *Codex) ExtractConfigServerURLs(data map[string]interface{}) []string {
	if provider, ok := data["provider"].(string); ok && provider != "" {
		return []string{provider}
	}
	return nil
}
