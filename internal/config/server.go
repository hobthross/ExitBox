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

package config

// ServerConfig holds the user-provided server configuration used when
// generating agent-specific configuration files.
type ServerConfig struct {
	ProviderID   string // slug, e.g. "local"
	ProviderName string // display name, e.g. "Qwen3.5-397B (local)"
	BaseURL      string // e.g. "http://localhost:8080/v1"
	APIKey       string // may be empty or vault reference
	ModelID      string // e.g. "qwen3.5-397b"
	ModelName    string // e.g. "Qwen3.5-397B-A17B"
	VaultKeyName string // non-empty when key is stored in vault
	Compaction   bool   // OpenCode: enable auto compaction with pruning
}

