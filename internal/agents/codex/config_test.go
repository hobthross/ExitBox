// ExitBox - Multi-Agent Container Sandbox
// Copyright (C) 2026 Cloud Exit B.V.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
