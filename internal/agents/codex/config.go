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
