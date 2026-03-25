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

package skills

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloud-exit/exitbox/internal/config"
	"gopkg.in/yaml.v3"
)

// Skill represents an installed skill.
type Skill struct {
	Name        string
	Description string
	Dir         string // absolute path to the skill directory
}

// WorkspaceSkillsDir returns the host path for a workspace's shared skills directory.
func WorkspaceSkillsDir(workspaceName string) string {
	return filepath.Join(config.Home, "profiles", "global", workspaceName, "skills")
}

// List returns all installed skills for a workspace.
func List(workspaceName string) ([]Skill, error) {
	dir := WorkspaceSkillsDir(workspaceName)
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var skills []Skill
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		skillFile := filepath.Join(dir, e.Name(), "SKILL.md")
		data, readErr := os.ReadFile(skillFile)
		if readErr != nil {
			continue
		}
		name, desc := parseFrontmatter(data)
		if name == "" {
			name = e.Name()
		}
		skills = append(skills, Skill{
			Name:        name,
			Description: desc,
			Dir:         filepath.Join(dir, e.Name()),
		})
	}
	return skills, nil
}

// Install writes a skill to the workspace skills directory.
// If additional files are provided, they are written alongside SKILL.md.
func Install(workspaceName, name string, files map[string][]byte) error {
	skillMd, ok := files["SKILL.md"]
	if !ok {
		return fmt.Errorf("no SKILL.md found")
	}

	dir := filepath.Join(WorkspaceSkillsDir(workspaceName), name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating skill directory: %w", err)
	}

	for relPath, content := range files {
		target := filepath.Join(dir, relPath)
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return fmt.Errorf("creating directory for %s: %w", relPath, err)
		}
		if err := os.WriteFile(target, content, 0644); err != nil {
			return fmt.Errorf("writing %s: %w", relPath, err)
		}
	}

	_ = skillMd // used via files map
	return nil
}

// Remove deletes a skill from the workspace.
func Remove(workspaceName, name string) error {
	dir := filepath.Join(WorkspaceSkillsDir(workspaceName), name)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("skill '%s' not found", name)
	}
	return os.RemoveAll(dir)
}

// Exists returns true if a skill is installed.
func Exists(workspaceName, name string) bool {
	skillFile := filepath.Join(WorkspaceSkillsDir(workspaceName), name, "SKILL.md")
	_, err := os.Stat(skillFile)
	return err == nil
}

// frontmatter is the YAML structure at the top of SKILL.md.
type frontmatter struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// parseFrontmatter extracts name and description from SKILL.md YAML frontmatter.
func parseFrontmatter(data []byte) (name, description string) {
	content := string(data)
	if !strings.HasPrefix(content, "---") {
		return "", ""
	}
	end := strings.Index(content[3:], "---")
	if end < 0 {
		return "", ""
	}
	var fm frontmatter
	if err := yaml.Unmarshal([]byte(content[3:3+end]), &fm); err != nil {
		return "", ""
	}
	return fm.Name, fm.Description
}
