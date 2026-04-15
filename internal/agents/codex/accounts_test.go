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
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestAccountManagerSaveAndSwitch(t *testing.T) {
	root := t.TempDir()
	manager := NewAccountManager(root)

	writeAccountFixture(t, filepath.Join(root, ".codex"), "auth.json", "personal-auth")
	writeAccountFixture(t, filepath.Join(root, ".config", "codex"), "config.json", "personal-config")

	if err := manager.Save("personal"); err != nil {
		t.Fatalf("Save(personal) error: %v", err)
	}

	state, err := manager.LoadState()
	if err != nil {
		t.Fatalf("LoadState() error: %v", err)
	}
	if state.Current != "personal" {
		t.Fatalf("Current = %q, want personal", state.Current)
	}

	writeAccountFixture(t, filepath.Join(root, ".codex"), "auth.json", "work-auth")
	writeAccountFixture(t, filepath.Join(root, ".config", "codex"), "config.json", "work-config")

	if err := manager.Save("work"); err != nil {
		t.Fatalf("Save(work) error: %v", err)
	}
	if err := manager.Switch("personal"); err != nil {
		t.Fatalf("Switch(personal) error: %v", err)
	}

	assertFileContent(t, filepath.Join(root, ".codex", "auth.json"), "personal-auth")
	assertFileContent(t, filepath.Join(root, ".config", "codex", "config.json"), "personal-config")
	assertFileContent(t, filepath.Join(root, ".exitbox", "accounts", "work", ".codex", "auth.json"), "work-auth")

	state, err = manager.LoadState()
	if err != nil {
		t.Fatalf("LoadState() after switch error: %v", err)
	}
	if state.Current != "personal" {
		t.Fatalf("Current after switch = %q, want personal", state.Current)
	}
	if state.Previous != "work" {
		t.Fatalf("Previous after switch = %q, want work", state.Previous)
	}
}

func TestAccountManagerAddClearsActiveConfig(t *testing.T) {
	root := t.TempDir()
	manager := NewAccountManager(root)

	writeAccountFixture(t, filepath.Join(root, ".codex"), "auth.json", "personal-auth")
	if err := manager.Save("personal"); err != nil {
		t.Fatalf("Save(personal) error: %v", err)
	}

	writeAccountFixture(t, filepath.Join(root, ".codex"), "auth.json", "transient-auth")
	if err := manager.Add("work"); err != nil {
		t.Fatalf("Add(work) error: %v", err)
	}

	if manager.ActiveDataExists() {
		t.Fatal("Add(work) should clear the active Codex config")
	}

	state, err := manager.LoadState()
	if err != nil {
		t.Fatalf("LoadState() error: %v", err)
	}
	if state.Current != "work" {
		t.Fatalf("Current = %q, want work", state.Current)
	}
	if state.Previous != "personal" {
		t.Fatalf("Previous = %q, want personal", state.Previous)
	}

	accounts, err := manager.List()
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(accounts) != 2 {
		t.Fatalf("List() returned %d accounts, want 2", len(accounts))
	}

	var foundWork bool
	for _, account := range accounts {
		if account.Name == "work" {
			foundWork = true
			if account.Ready {
				t.Fatal("newly added account should be pending login")
			}
		}
	}
	if !foundWork {
		t.Fatal("expected to find work account")
	}
}

func TestAccountManagerSwitchRequiresNamedCurrentAccount(t *testing.T) {
	root := t.TempDir()
	manager := NewAccountManager(root)

	writeAccountFixture(t, filepath.Join(root, ".codex"), "auth.json", "unnamed-auth")
	writeAccountFixture(t, filepath.Join(root, ".exitbox", "accounts", "personal", ".codex"), "auth.json", "personal-auth")

	err := manager.Switch("personal")
	if !errors.Is(err, ErrUnnamedActiveAccount) {
		t.Fatalf("Switch(personal) error = %v, want %v", err, ErrUnnamedActiveAccount)
	}
}

func TestValidateAccountName(t *testing.T) {
	cases := []struct {
		name    string
		wantErr bool
	}{
		{name: "personal"},
		{name: "work-2"},
		{name: "", wantErr: true},
		{name: ".", wantErr: true},
		{name: "foo/bar", wantErr: true},
		{name: `foo\bar`, wantErr: true},
	}

	for _, tc := range cases {
		err := validateAccountName(tc.name)
		if tc.wantErr && err == nil {
			t.Fatalf("validateAccountName(%q) expected error", tc.name)
		}
		if !tc.wantErr && err != nil {
			t.Fatalf("validateAccountName(%q) unexpected error: %v", tc.name, err)
		}
	}
}

func writeAccountFixture(t *testing.T, dir, name, value string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("MkdirAll(%s): %v", dir, err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte(value), 0644); err != nil {
		t.Fatalf("WriteFile(%s): %v", filepath.Join(dir, name), err)
	}
}

func assertFileContent(t *testing.T, path, want string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%s): %v", path, err)
	}
	if string(data) != want {
		t.Fatalf("%s = %q, want %q", path, string(data), want)
	}
}
