package opencode

import (
	"path/filepath"

	"github.com/cloud-exit/exitbox/internal/config"
)

// GenerateConfig produces an OpenCode config map.
func (o *OpenCode) GenerateConfig(cfg config.ServerConfig) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"$schema": "https://opencode.ai/config.json",
		"provider": map[string]interface{}{
			cfg.ProviderID: map[string]interface{}{
				"npm":  "@ai-sdk/openai-compatible",
				"name": cfg.ProviderName,
				"options": map[string]interface{}{
					"baseURL": cfg.BaseURL,
				},
				"models": map[string]interface{}{
					cfg.ModelID: map[string]interface{}{
						"name": cfg.ModelName,
					},
				},
			},
		},
		"model": cfg.ProviderID + "/" + cfg.ModelID,
	}
	if cfg.Compaction {
		result["compaction"] = map[string]interface{}{
			"auto":  true,
			"prune": true,
		}
	}
	return result, nil
}

// LogSearchDirs returns directories to search for OpenCode log files.
func (o *OpenCode) LogSearchDirs(home, agentCfgDir string) []string {
	return []string{
		filepath.Join(home, ".local", "share", "opencode", "log"),
		filepath.Join(home, ".local", "share", "opencode", "logs"),
		filepath.Join(home, ".opencode"),
		filepath.Join(agentCfgDir, ".opencode"),
		filepath.Join(agentCfgDir, ".config", "opencode"),
	}
}

func (o *OpenCode) ConfigFilePath(agentDir string) string {
	return filepath.Join(agentDir, ".config", "opencode", "opencode.json")
}

// ExtractConfigServerURLs walks provider.*.options.baseURL in an OpenCode config.
func (o *OpenCode) ExtractConfigServerURLs(data map[string]interface{}) []string {
	providers, ok := data["provider"].(map[string]interface{})
	if !ok {
		return nil
	}
	var urls []string
	for _, pv := range providers {
		provider, ok := pv.(map[string]interface{})
		if !ok {
			continue
		}
		opts, ok := provider["options"].(map[string]interface{})
		if !ok {
			continue
		}
		if baseURL, ok := opts["baseURL"].(string); ok && baseURL != "" {
			urls = append(urls, baseURL)
		}
	}
	return urls
}
