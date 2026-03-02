package claude

import (
	"github.com/cloud-exit/exitbox/internal/generate"
)

// GenerateConfig produces a Claude Code settings.json config map.
func (c *Claude) GenerateConfig(cfg generate.ServerConfig) (map[string]interface{}, error) {
	m := map[string]interface{}{
		"apiBaseUrl": cfg.BaseURL,
		"model":      cfg.ModelID,
	}
	if cfg.APIKey != "" && cfg.VaultKeyName == "" {
		m["apiKey"] = cfg.APIKey
	}
	return m, nil
}
