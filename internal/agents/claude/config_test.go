package claude

import (
	"testing"

	"github.com/cloud-exit/exitbox/internal/config"
)

func TestGenerateConfig_Claude(t *testing.T) {
	t.Run("with api key no vault", func(t *testing.T) {
		cfg := config.ServerConfig{
			BaseURL: "http://localhost:8080/v1",
			ModelID: "gpt-4",
			APIKey:  "sk-test",
		}

		c := &Claude{}
		result, err := c.GenerateConfig(cfg)
		if err != nil {
			t.Fatalf("GenerateConfig error: %v", err)
		}

		if result["apiBaseUrl"] != "http://localhost:8080/v1" {
			t.Errorf("wrong apiBaseUrl: %v", result["apiBaseUrl"])
		}
		if result["model"] != "gpt-4" {
			t.Errorf("wrong model: %v", result["model"])
		}
		if result["apiKey"] != "sk-test" {
			t.Errorf("wrong apiKey: %v", result["apiKey"])
		}
	})

	t.Run("with vault key", func(t *testing.T) {
		cfg := config.ServerConfig{
			BaseURL:      "http://localhost:8080/v1",
			ModelID:      "gpt-4",
			APIKey:       "sk-test",
			VaultKeyName: "API_KEY",
		}

		c := &Claude{}
		result, err := c.GenerateConfig(cfg)
		if err != nil {
			t.Fatalf("GenerateConfig error: %v", err)
		}

		if _, ok := result["apiKey"]; ok {
			t.Error("apiKey should not be set when vault key is used")
		}
	})

	t.Run("no api key", func(t *testing.T) {
		cfg := config.ServerConfig{
			BaseURL: "http://localhost:8080/v1",
			ModelID: "gpt-4",
		}

		c := &Claude{}
		result, err := c.GenerateConfig(cfg)
		if err != nil {
			t.Fatalf("GenerateConfig error: %v", err)
		}

		if _, ok := result["apiKey"]; ok {
			t.Error("apiKey should not be set when empty")
		}
	})
}

