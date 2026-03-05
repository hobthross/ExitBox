package qwen

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/cloud-exit/exitbox/internal/config"
)

func TestGenerateConfig_Qwen(t *testing.T) {
	cfg := config.ServerConfig{
		ProviderID:   "local",
		ProviderName: "Local Server",
		BaseURL:      "http://localhost:8080/v1",
		ModelID:      "qwen3-coder-plus",
		ModelName:    "Qwen3 Coder Plus",
	}

	q := &Qwen{}
	result, err := q.GenerateConfig(cfg)
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

	providers, ok := parsed["modelProviders"].(map[string]interface{})
	if !ok {
		t.Fatal("missing modelProviders")
	}
	openaiList, ok := providers["openai"].([]interface{})
	if !ok || len(openaiList) != 1 {
		t.Fatal("modelProviders.openai should be a single-element array")
	}
	provider, ok := openaiList[0].(map[string]interface{})
	if !ok {
		t.Fatal("first openai provider should be an object")
	}
	if provider["baseUrl"] != "http://localhost:8080/v1" {
		t.Errorf("expected baseUrl, got %v", provider["baseUrl"])
	}
	if provider["id"] != "qwen3-coder-plus" {
		t.Errorf("expected id qwen3-coder-plus, got %v", provider["id"])
	}

	model, ok := parsed["model"].(map[string]interface{})
	if !ok {
		t.Fatal("missing model")
	}
	if model["name"] != "qwen3-coder-plus" {
		t.Errorf("expected model.name qwen3-coder-plus, got %v", model["name"])
	}

	security, ok := parsed["security"].(map[string]interface{})
	if !ok {
		t.Fatal("missing security")
	}
	auth, ok := security["auth"].(map[string]interface{})
	if !ok {
		t.Fatal("missing security.auth")
	}
	if auth["selectedType"] != "openai" {
		t.Errorf("expected selectedType openai, got %v", auth["selectedType"])
	}
}

func TestGenerateConfig_WithAPIKey(t *testing.T) {
	cfg := config.ServerConfig{
		ProviderID:   "local",
		ProviderName: "Local",
		BaseURL:      "http://localhost:8080/v1",
		ModelID:      "qwen",
		ModelName:    "Qwen",
		APIKey:       "sk-test",
	}

	q := &Qwen{}
	result, err := q.GenerateConfig(cfg)
	if err != nil {
		t.Fatalf("GenerateConfig error: %v", err)
	}
	env, ok := result["env"].(map[string]interface{})
	if !ok {
		t.Fatal("expected env when APIKey is set")
	}
	if env["OPENAI_API_KEY"] != "sk-test" {
		t.Errorf("expected OPENAI_API_KEY in env, got %v", env["OPENAI_API_KEY"])
	}
}

func TestConfigFilePath_Qwen(t *testing.T) {
	q := &Qwen{}
	got := q.ConfigFilePath("/base")
	want := filepath.Join("/base", ".config", "qwen", "settings.json")
	if got != want {
		t.Errorf("ConfigFilePath(/base) = %q, want %q", got, want)
	}
}

func TestExtractConfigServerURLs_Qwen(t *testing.T) {
	q := &Qwen{}

	t.Run("empty", func(t *testing.T) {
		urls := q.ExtractConfigServerURLs(map[string]interface{}{})
		if len(urls) != 0 {
			t.Errorf("expected no URLs, got %v", urls)
		}
	})

	t.Run("openai baseUrl", func(t *testing.T) {
		data := map[string]interface{}{
			"modelProviders": map[string]interface{}{
				"openai": []interface{}{
					map[string]interface{}{
						"id":     "m1",
						"baseUrl": "https://api.example.com/v1",
					},
				},
			},
		}
		urls := q.ExtractConfigServerURLs(data)
		if len(urls) != 1 || urls[0] != "https://api.example.com/v1" {
			t.Errorf("expected one URL https://api.example.com/v1, got %v", urls)
		}
	})
}
