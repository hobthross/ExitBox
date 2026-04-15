package plugins

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestProjectPluginsDir_Codex(t *testing.T) {
	projectDir := t.TempDir()
	got, err := ProjectPluginsDir(projectDir, "codex")
	if err != nil {
		t.Fatalf("ProjectPluginsDir returned error: %v", err)
	}
	want := filepath.Join(projectDir, "plugins")
	if got != want {
		t.Fatalf("ProjectPluginsDir = %q, want %q", got, want)
	}
}

func TestProjectPluginsDir_UnsupportedAgent(t *testing.T) {
	if _, err := ProjectPluginsDir(t.TempDir(), "claude"); err == nil {
		t.Fatal("ProjectPluginsDir should reject unsupported agents")
	}
}

func TestInstallListAndRemove(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	projectDir := t.TempDir()
	repo := initMarketplaceRepo(t, "caveman-repo")

	installed, err := Install(projectDir, "codex", repo, "")
	if err != nil {
		t.Fatalf("Install returned error: %v", err)
	}
	if len(installed) != 1 || installed[0].Name != "caveman" {
		t.Fatalf("Install = %#v, want caveman", installed)
	}
	if _, err := os.Stat(filepath.Join(projectDir, "plugins", "caveman", ".codex-plugin", "plugin.json")); err != nil {
		t.Fatalf("expected plugin payload in project plugins dir: %v", err)
	}

	marketplacePath := filepath.Join(projectDir, ".agents", "plugins", "marketplace.json")
	data, err := os.ReadFile(marketplacePath)
	if err != nil {
		t.Fatalf("expected marketplace file: %v", err)
	}
	var m struct {
		Plugins []struct {
			Name   string `json:"name"`
			Source struct {
				Path string `json:"path"`
			} `json:"source"`
		} `json:"plugins"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("invalid marketplace json: %v", err)
	}
	if len(m.Plugins) != 1 || m.Plugins[0].Name != "caveman" || m.Plugins[0].Source.Path != "./plugins/caveman" {
		t.Fatalf("unexpected marketplace content: %#v", m.Plugins)
	}

	listed, err := List(projectDir, "codex")
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(listed) != 1 || listed[0].Name != "caveman" {
		t.Fatalf("List = %#v, want caveman", listed)
	}

	if err := Remove(projectDir, "codex", "caveman"); err != nil {
		t.Fatalf("Remove returned error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(projectDir, "plugins", "caveman")); !os.IsNotExist(err) {
		t.Fatalf("plugin dir still exists after remove: %v", err)
	}
}

func initMarketplaceRepo(t *testing.T, name string) string {
	t.Helper()

	repo := filepath.Join(t.TempDir(), name)
	marketplacePath := filepath.Join(repo, ".agents", "plugins")
	pluginPath := filepath.Join(repo, "plugins", "caveman", ".codex-plugin")

	if err := os.MkdirAll(marketplacePath, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(pluginPath, 0755); err != nil {
		t.Fatal(err)
	}

	marketplaceJSON := `{
  "name": "caveman-repo",
  "interface": {
    "displayName": "Caveman Repo"
  },
  "plugins": [
    {
      "name": "caveman",
      "source": {
        "source": "local",
        "path": "./plugins/caveman"
      },
      "policy": {
        "installation": "AVAILABLE",
        "authentication": "ON_INSTALL"
      },
      "category": "Productivity"
    }
  ]
}`
	if err := os.WriteFile(filepath.Join(marketplacePath, "marketplace.json"), []byte(marketplaceJSON), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pluginPath, "plugin.json"), []byte(`{"name":"caveman"}`), 0644); err != nil {
		t.Fatal(err)
	}

	runGit(t, repo, "init")
	runGit(t, repo, "config", "user.email", "test@example.com")
	runGit(t, repo, "config", "user.name", "Test User")
	runGit(t, repo, "add", ".")
	runGit(t, repo, "commit", "-m", "init")
	return repo
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, string(out))
	}
}
