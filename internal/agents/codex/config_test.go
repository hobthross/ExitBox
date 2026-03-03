package codex

import (
	"path/filepath"
	"testing"

	"github.com/cloud-exit/exitbox/internal/config"
)

func TestGenerateConfig_Codex(t *testing.T) {
	cfg := config.ServerConfig{
		ProviderID: "local",
		BaseURL:    "http://localhost:8080/v1",
		ModelID:    "gpt-4",
	}

	c := &Codex{}
	result, err := c.GenerateConfig(cfg)
	if err != nil {
		t.Fatalf("GenerateConfig error: %v", err)
	}

	if result["model"] != "local/gpt-4" {
		t.Errorf("wrong model: %v", result["model"])
	}
	if result["provider"] != "http://localhost:8080/v1" {
		t.Errorf("wrong provider: %v", result["provider"])
	}
}

func TestConfigFilePath_Codex(t *testing.T) {
	c := &Codex{}
	got := c.ConfigFilePath("/base")
	want := filepath.Join("/base", ".codex", "config.json")
	if got != want {
		t.Errorf("ConfigFilePath(/base) = %q, want %q", got, want)
	}
}

