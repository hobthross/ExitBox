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
	"path/filepath"

	"github.com/cloud-exit/exitbox/internal/config"
)

// GenerateConfig produces a Qwen Code settings.json config map (OpenAI-compatible provider).
func (q *Qwen) GenerateConfig(cfg config.ServerConfig) (map[string]interface{}, error) {
	provider := map[string]interface{}{
		"id":          cfg.ModelID,
		"name":        cfg.ModelName,
		"baseUrl":     cfg.BaseURL,
		"description": cfg.ProviderName,
		"envKey":      "OPENAI_API_KEY",
		"models": map[string]interface{}{
			cfg.ModelID: map[string]interface{}{
				"name": cfg.ModelName,
			},
		},
	}
	result := map[string]interface{}{
		"modelProviders": map[string]interface{}{
			"openai": []map[string]interface{}{provider},
		},
		"security": map[string]interface{}{
			"auth": map[string]interface{}{
				"selectedType": "openai",
			},
		},
		"model": map[string]interface{}{
			"name": cfg.ModelID,
		},
	}
	if cfg.APIKey != "" && cfg.VaultKeyName == "" {
		result["env"] = map[string]interface{}{
			"OPENAI_API_KEY": cfg.APIKey,
		}
	}
	return result, nil
}

// LogSearchDirs returns directories to search for Qwen Code log files.
func (q *Qwen) LogSearchDirs(home, agentCfgDir string) []string {
	return []string{
		filepath.Join(home, ".qwen"),
		filepath.Join(home, ".config", "qwen"),
		filepath.Join(agentCfgDir, ".qwen"),
		filepath.Join(agentCfgDir, ".config", "qwen"),
	}
}

func (q *Qwen) ConfigFilePath(agentDir string) string {
	return filepath.Join(agentDir, ".config", "qwen", "settings.json")
}

// ExtractConfigServerURLs walks modelProviders.*[].baseUrl in Qwen settings.json.
func (q *Qwen) ExtractConfigServerURLs(data map[string]interface{}) []string {
	providers, ok := data["modelProviders"].(map[string]interface{})
	if !ok {
		return nil
	}
	var urls []string
	for _, pv := range providers {
		list, ok := pv.([]interface{})
		if !ok {
			continue
		}
		for _, item := range list {
			entry, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if baseURL, ok := entry["baseUrl"].(string); ok && baseURL != "" {
				urls = append(urls, baseURL)
			}
		}
	}
	return urls
}
