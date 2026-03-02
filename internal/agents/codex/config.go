package codex

import (
	"github.com/cloud-exit/exitbox/internal/generate"
)

// GenerateConfig produces a Codex config.json config map.
func (c *Codex) GenerateConfig(cfg generate.ServerConfig) (map[string]interface{}, error) {
	return map[string]interface{}{
		"model":    cfg.ProviderID + "/" + cfg.ModelID,
		"provider": cfg.BaseURL,
	}, nil
}
