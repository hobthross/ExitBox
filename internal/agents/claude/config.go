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

package claude

import (
	"path/filepath"

	"github.com/cloud-exit/exitbox/internal/config"
)

// GenerateConfig produces a Claude Code settings.json config map.
func (c *Claude) GenerateConfig(cfg config.ServerConfig) (map[string]interface{}, error) {
	m := map[string]interface{}{
		"apiBaseUrl": cfg.BaseURL,
		"model":      cfg.ModelID,
	}
	if cfg.APIKey != "" && cfg.VaultKeyName == "" {
		m["apiKey"] = cfg.APIKey
	}
	return m, nil
}

// LogSearchDirs returns directories to search for Claude Code log files.
func (c *Claude) LogSearchDirs(home, agentCfgDir string) []string {
	return []string{
		filepath.Join(home, ".claude"),
		filepath.Join(agentCfgDir, ".claude"),
	}
}

func (c *Claude) ConfigFilePath(agentDir string) string {
	return filepath.Join(agentDir, ".claude", "settings.json")
}

// ExtractConfigServerURLs reads the top-level apiBaseUrl from a Claude config.
func (c *Claude) ExtractConfigServerURLs(data map[string]interface{}) []string {
	if baseURL, ok := data["apiBaseUrl"].(string); ok && baseURL != "" {
		return []string{baseURL}
	}
	return nil
}
