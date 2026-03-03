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
