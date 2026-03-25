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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var githubTreeRE = regexp.MustCompile(`^https?://github\.com/([^/]+)/([^/]+)/tree/([^/]+)/(.+)$`)

// SourceType identifies the type of skill source.
type SourceType int

const (
	SourceUnknown    SourceType = iota
	SourceGitHubTree            // GitHub tree URL (directory)
	SourceRawURL                // Direct URL to a file
	SourceLocalPath             // Local filesystem path
)

// DetectSource classifies a source string.
func DetectSource(input string) SourceType {
	if githubTreeRE.MatchString(input) {
		return SourceGitHubTree
	}
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		return SourceRawURL
	}
	return SourceLocalPath
}

// FetchResult contains the fetched skill files and derived metadata.
type FetchResult struct {
	Name  string            // derived skill name
	Files map[string][]byte // relative path → content (always includes "SKILL.md")
}

// Fetch retrieves skill file(s) from the given source.
func Fetch(source string) (*FetchResult, error) {
	switch DetectSource(source) {
	case SourceGitHubTree:
		return fetchGitHubTree(source)
	case SourceRawURL:
		return fetchRawURL(source)
	case SourceLocalPath:
		return fetchLocal(source)
	default:
		return nil, fmt.Errorf("unsupported source: %s", source)
	}
}

// fetchGitHubTree handles GitHub tree URLs like:
// https://github.com/anthropics/skills/tree/main/skills/frontend-design
func fetchGitHubTree(url string) (*FetchResult, error) {
	m := githubTreeRE.FindStringSubmatch(url)
	if m == nil {
		return nil, fmt.Errorf("invalid GitHub tree URL: %s", url)
	}
	owner, repo, ref, dirPath := m[1], m[2], m[3], m[4]

	// Recursively fetch all files including nested directories.
	files := make(map[string][]byte)
	if err := githubFetchDir(owner, repo, ref, dirPath, "", files); err != nil {
		return nil, fmt.Errorf("fetching GitHub directory: %w", err)
	}

	if _, ok := files["SKILL.md"]; !ok {
		return nil, fmt.Errorf("no SKILL.md found in %s", url)
	}

	// Derive name from the last path component or frontmatter.
	name := filepath.Base(dirPath)
	if fmName, _ := parseFrontmatter(files["SKILL.md"]); fmName != "" {
		name = fmName
	}

	return &FetchResult{Name: name, Files: files}, nil
}

// githubFetchDir recursively fetches all files from a GitHub directory.
// prefix is the relative path prefix for nested directories.
func githubFetchDir(owner, repo, ref, dirPath, prefix string, files map[string][]byte) error {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s", owner, repo, dirPath, ref)
	entries, err := githubListDir(apiURL)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		relPath := entry.Name
		if prefix != "" {
			relPath = prefix + "/" + entry.Name
		}

		switch entry.Type {
		case "file":
			if entry.DownloadURL == "" {
				continue
			}
			content, dlErr := httpGet(entry.DownloadURL)
			if dlErr != nil {
				continue
			}
			files[relPath] = content
		case "dir":
			subPath := dirPath + "/" + entry.Name
			if err := githubFetchDir(owner, repo, ref, subPath, relPath, files); err != nil {
				continue // best-effort for subdirectories
			}
		}
	}
	return nil
}

// fetchRawURL handles direct URLs to a SKILL.md file.
func fetchRawURL(url string) (*FetchResult, error) {
	content, err := httpGet(url)
	if err != nil {
		return nil, fmt.Errorf("downloading %s: %w", url, err)
	}

	// Derive name from URL path or frontmatter.
	name := ""
	if fmName, _ := parseFrontmatter(content); fmName != "" {
		name = fmName
	}
	if name == "" {
		// Use parent directory from URL path.
		parts := strings.Split(strings.TrimSuffix(url, "/"), "/")
		for i := len(parts) - 1; i >= 0; i-- {
			if !strings.EqualFold(parts[i], "SKILL.md") && parts[i] != "" {
				name = parts[i]
				break
			}
		}
	}
	if name == "" {
		name = "unnamed-skill"
	}

	return &FetchResult{
		Name:  name,
		Files: map[string][]byte{"SKILL.md": content},
	}, nil
}

// fetchLocal reads a skill from a local path (file or directory).
func fetchLocal(path string) (*FetchResult, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	if !info.IsDir() {
		// Single file — treat as SKILL.md.
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil, readErr
		}
		name := ""
		if fmName, _ := parseFrontmatter(content); fmName != "" {
			name = fmName
		}
		if name == "" {
			name = filepath.Base(filepath.Dir(path))
		}
		return &FetchResult{
			Name:  name,
			Files: map[string][]byte{"SKILL.md": content},
		}, nil
	}

	// Directory — read all files, require SKILL.md.
	files := make(map[string][]byte)
	err = filepath.WalkDir(path, func(p string, d os.DirEntry, walkErr error) error {
		if walkErr != nil || d.IsDir() {
			return walkErr
		}
		rel, relErr := filepath.Rel(path, p)
		if relErr != nil {
			return relErr
		}
		content, readErr := os.ReadFile(p)
		if readErr != nil {
			return readErr
		}
		files[rel] = content
		return nil
	})
	if err != nil {
		return nil, err
	}

	if _, ok := files["SKILL.md"]; !ok {
		return nil, fmt.Errorf("no SKILL.md found in %s", path)
	}

	name := filepath.Base(path)
	if fmName, _ := parseFrontmatter(files["SKILL.md"]); fmName != "" {
		name = fmName
	}

	return &FetchResult{Name: name, Files: files}, nil
}

// githubEntry represents a file entry from the GitHub Contents API.
type githubEntry struct {
	Name        string `json:"name"`
	Type        string `json:"type"` // "file" or "dir"
	DownloadURL string `json:"download_url"`
}

func githubListDir(apiURL string) ([]githubEntry, error) {
	data, err := httpGet(apiURL)
	if err != nil {
		return nil, err
	}
	var entries []githubEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("parsing GitHub API response: %w", err)
	}
	return entries, nil
}

func httpGet(url string) ([]byte, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d for %s", resp.StatusCode, url)
	}
	return io.ReadAll(resp.Body)
}
