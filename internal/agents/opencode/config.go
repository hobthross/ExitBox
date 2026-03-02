package opencode

import (
	"github.com/cloud-exit/exitbox/internal/generate"
)

// GenerateConfig produces an OpenCode config map.
func (o *OpenCode) GenerateConfig(cfg generate.ServerConfig) (map[string]interface{}, error) {
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
