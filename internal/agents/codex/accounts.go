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
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cloud-exit/exitbox/internal/fsutil"
)

var (
	ErrNoActiveAccountData  = errors.New("no active Codex account data found")
	ErrUnnamedActiveAccount = errors.New("current Codex account is unnamed; run 'exitbox codex accounts save <name>' first")
)

type AccountState struct {
	Current  string `json:"current,omitempty"`
	Previous string `json:"previous,omitempty"`
}

type AccountInfo struct {
	Name     string
	Current  bool
	Previous bool
	Ready    bool
}

type AccountManager struct {
	Root string
}

func NewAccountManager(root string) *AccountManager {
	return &AccountManager{Root: root}
}

func (m *AccountManager) List() ([]AccountInfo, error) {
	state, err := m.LoadState()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(m.accountsRoot())
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	items := make([]AccountInfo, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		items = append(items, AccountInfo{
			Name:     name,
			Current:  name == state.Current,
			Previous: name == state.Previous,
			Ready:    m.accountHasData(name),
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
	})

	return items, nil
}

func (m *AccountManager) LoadState() (AccountState, error) {
	data, err := os.ReadFile(m.stateFile())
	if err != nil {
		if os.IsNotExist(err) {
			return AccountState{}, nil
		}
		return AccountState{}, err
	}

	var state AccountState
	if err := json.Unmarshal(data, &state); err != nil {
		return AccountState{}, err
	}
	return state, nil
}

func (m *AccountManager) Save(name string) error {
	if err := validateAccountName(name); err != nil {
		return err
	}
	if !m.ActiveDataExists() {
		return ErrNoActiveAccountData
	}
	if err := m.copyActiveToAccount(name); err != nil {
		return err
	}

	state, err := m.LoadState()
	if err != nil {
		return err
	}
	if state.Current != name {
		state.Previous = state.Current
	}
	state.Current = name
	return m.saveState(state)
}

func (m *AccountManager) Add(name string) error {
	if err := validateAccountName(name); err != nil {
		return err
	}
	if m.accountExists(name) {
		return fmt.Errorf("Codex account %q already exists", name)
	}

	state, err := m.LoadState()
	if err != nil {
		return err
	}
	if m.ActiveDataExists() {
		if state.Current == "" {
			return ErrUnnamedActiveAccount
		}
		if err := m.copyActiveToAccount(state.Current); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(m.accountDir(name), 0755); err != nil {
		return err
	}
	if err := m.clearActive(); err != nil {
		return err
	}

	state.Previous = state.Current
	state.Current = name
	return m.saveState(state)
}

func (m *AccountManager) Switch(name string) error {
	if err := validateAccountName(name); err != nil {
		return err
	}

	state, err := m.LoadState()
	if err != nil {
		return err
	}
	if state.Current == name {
		return nil
	}
	if !m.accountHasData(name) {
		return fmt.Errorf("Codex account %q has no saved login yet", name)
	}

	if m.ActiveDataExists() {
		if state.Current == "" {
			return ErrUnnamedActiveAccount
		}
		if err := m.copyActiveToAccount(state.Current); err != nil {
			return err
		}
	}

	if err := m.clearActive(); err != nil {
		return err
	}
	if err := m.copyAccountToActive(name); err != nil {
		return err
	}

	state.Previous = state.Current
	state.Current = name
	return m.saveState(state)
}

func (m *AccountManager) ActiveDataExists() bool {
	return dirHasEntries(m.activeCodexDir()) || dirHasEntries(m.activeConfigDir())
}

func (m *AccountManager) accountExists(name string) bool {
	info, err := os.Stat(m.accountDir(name))
	return err == nil && info.IsDir()
}

func (m *AccountManager) accountHasData(name string) bool {
	return dirHasEntries(filepath.Join(m.accountDir(name), ".codex")) ||
		dirHasEntries(filepath.Join(m.accountDir(name), ".config", "codex"))
}

func (m *AccountManager) copyActiveToAccount(name string) error {
	if !m.ActiveDataExists() {
		return ErrNoActiveAccountData
	}

	root := m.accountDir(name)
	if err := os.RemoveAll(root); err != nil {
		return err
	}
	if err := os.MkdirAll(root, 0755); err != nil {
		return err
	}
	if err := copyIfPresent(m.activeCodexDir(), filepath.Join(root, ".codex")); err != nil {
		return err
	}
	if err := copyIfPresent(m.activeConfigDir(), filepath.Join(root, ".config", "codex")); err != nil {
		return err
	}
	return nil
}

func (m *AccountManager) copyAccountToActive(name string) error {
	root := m.accountDir(name)
	if err := copyIfPresent(filepath.Join(root, ".codex"), m.activeCodexDir()); err != nil {
		return err
	}
	if err := copyIfPresent(filepath.Join(root, ".config", "codex"), m.activeConfigDir()); err != nil {
		return err
	}
	return nil
}

func (m *AccountManager) clearActive() error {
	if err := os.RemoveAll(m.activeCodexDir()); err != nil {
		return err
	}
	if err := os.RemoveAll(m.activeConfigDir()); err != nil {
		return err
	}
	return nil
}

func (m *AccountManager) saveState(state AccountState) error {
	if err := os.MkdirAll(filepath.Dir(m.stateFile()), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(m.stateFile(), data, 0644)
}

func (m *AccountManager) accountsRoot() string {
	return filepath.Join(m.Root, ".exitbox", "accounts")
}

func (m *AccountManager) accountDir(name string) string {
	return filepath.Join(m.accountsRoot(), name)
}

func (m *AccountManager) stateFile() string {
	return filepath.Join(m.accountsRoot(), "state.json")
}

func (m *AccountManager) activeCodexDir() string {
	return filepath.Join(m.Root, ".codex")
}

func (m *AccountManager) activeConfigDir() string {
	return filepath.Join(m.Root, ".config", "codex")
}

func validateAccountName(name string) error {
	name = strings.TrimSpace(name)
	switch {
	case name == "":
		return errors.New("Codex account name cannot be empty")
	case name == "." || name == "..":
		return fmt.Errorf("invalid Codex account name %q", name)
	case strings.ContainsAny(name, `/\`):
		return fmt.Errorf("Codex account name %q cannot contain path separators", name)
	default:
		return nil
	}
}

func copyIfPresent(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", src)
	}
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}
	return fsutil.CopyDirRecursive(src, dst)
}

func dirHasEntries(path string) bool {
	entries, err := os.ReadDir(path)
	return err == nil && len(entries) > 0
}
