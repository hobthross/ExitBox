package opencode

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/cloud-exit/exitbox/internal/config"
)

func TestGenerateConfig_OpenCode(t *testing.T) {
	t.Run("without compaction", func(t *testing.T) {
		cfg := config.ServerConfig{
			ProviderID:   "local",
			ProviderName: "Local Server",
			BaseURL:      "http://localhost:8080/v1",
			ModelID:      "qwen3",
			ModelName:    "Qwen3",
		}

		o := &OpenCode{}
		result, err := o.GenerateConfig(cfg)
		if err != nil {
			t.Fatalf("GenerateConfig error: %v", err)
		}

		data, err := json.Marshal(result)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var parsed map[string]interface{}
		if err := json.Unmarshal(data, &parsed); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		if parsed["$schema"] != "https://opencode.ai/config.json" {
			t.Error("missing or wrong $schema")
		}

		providers, ok := parsed["provider"].(map[string]interface{})
		if !ok {
			t.Fatal("missing provider")
		}
		local, ok := providers["local"].(map[string]interface{})
		if !ok {
			t.Fatal("missing local provider")
		}
		if local["name"] != "Local Server" {
			t.Errorf("expected Local Server, got %v", local["name"])
		}

		if parsed["model"] != "local/qwen3" {
			t.Errorf("expected local/qwen3, got %v", parsed["model"])
		}

		if _, ok := parsed["compaction"]; ok {
			t.Error("compaction should not be present when disabled")
		}
	})

	t.Run("with compaction", func(t *testing.T) {
		cfg := config.ServerConfig{
			ProviderID:   "local",
			ProviderName: "Local Server",
			BaseURL:      "http://localhost:8080/v1",
			ModelID:      "qwen3",
			ModelName:    "Qwen3",
			Compaction:   true,
		}

		o := &OpenCode{}
		result, err := o.GenerateConfig(cfg)
		if err != nil {
			t.Fatalf("GenerateConfig error: %v", err)
		}

		data, err := json.Marshal(result)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var parsed map[string]interface{}
		if err := json.Unmarshal(data, &parsed); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		comp, ok := parsed["compaction"].(map[string]interface{})
		if !ok {
			t.Fatal("missing compaction")
		}
		if comp["auto"] != true {
			t.Errorf("expected compaction.auto=true, got %v", comp["auto"])
		}
		if comp["prune"] != true {
			t.Errorf("expected compaction.prune=true, got %v", comp["prune"])
		}
	})
}

func TestConfigFilePath_OpenCode(t *testing.T) {
	o := &OpenCode{}
	got := o.ConfigFilePath("/base")
	want := filepath.Join("/base", ".config", "opencode", "opencode.json")
	if got != want {
		t.Errorf("ConfigFilePath(/base) = %q, want %q", got, want)
	}
}

