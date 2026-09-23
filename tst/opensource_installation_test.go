package tst

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestOpenSourceInstallerIsIdempotent(t *testing.T) {
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	home := t.TempDir()
	bin := filepath.Join(t.TempDir(), "bin")
	if err := os.MkdirAll(bin, 0o755); err != nil {
		t.Fatal(err)
	}

	writeExecutable(t, filepath.Join(bin, "curl"), "#!/bin/sh\nprintf '{}\\n'\n")
	writeExecutable(t, filepath.Join(bin, "ollama"), `#!/bin/sh
case "$1" in
  --version) echo "ollama version 0.test" ;;
  list) printf 'NAME ID SIZE MODIFIED\nqwen3.5:9b abc 6.6GB now\n' ;;
  pull) echo "unexpected pull" >&2; exit 9 ;;
  *) exit 0 ;;
esac
`)
	writeExecutable(t, filepath.Join(bin, "qwen"), `#!/bin/sh
if [ "$1" = "--version" ]; then
  echo "0.test"
else
  printf '[{"type":"result","subtype":"success","result":"RAPIDOU_OK"}]\n'
fi
`)
	writeExecutable(t, filepath.Join(bin, "codium"), "#!/bin/sh\n[ \"$1\" = \"--list-extensions\" ] && echo qwenlm.qwen-code-vscode-ide-companion\n")

	settingsPath := filepath.Join(home, ".qwen", "settings.json")
	if err := os.MkdirAll(filepath.Dir(settingsPath), 0o755); err != nil {
		t.Fatal(err)
	}
	original := `{"theme":"Solarized","modelProviders":{"openai":[{"id":"cloud-model","baseUrl":"https://example.test/v1"}]}}` + "\n"
	if err := os.WriteFile(settingsPath, []byte(original), 0o600); err != nil {
		t.Fatal(err)
	}

	installer := filepath.Join(root, "run", "install-opensource.sh")
	environment := append(os.Environ(),
		"HOME="+home,
		"PATH="+bin+":/usr/bin:/bin",
		"RAPIDOU_QWEN_MODEL=qwen3.5:9b",
	)
	for attempt := 0; attempt < 2; attempt++ {
		command := exec.Command(installer)
		command.Env = environment
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("open-source install attempt %d failed: %v: %s", attempt+1, err, output)
		}
	}

	backup, err := os.ReadFile(settingsPath + ".rapidou-backup")
	if err != nil {
		t.Fatalf("missing settings backup: %v", err)
	}
	if string(backup) != original {
		t.Fatalf("backup changed: got %q, want %q", backup, original)
	}

	contents, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatal(err)
	}
	var settings map[string]any
	if err := json.Unmarshal(contents, &settings); err != nil {
		t.Fatalf("invalid merged settings: %v", err)
	}
	if settings["theme"] != "Solarized" {
		t.Fatalf("unrelated setting was not preserved: %s", contents)
	}
	providers := settings["modelProviders"].(map[string]any)["openai"].([]any)
	rapidouEntries := 0
	for _, raw := range providers {
		entry := raw.(map[string]any)
		if entry["description"] == "Rapidou local Ollama model" {
			rapidouEntries++
		}
	}
	if len(providers) != 2 || rapidouEntries != 1 {
		t.Fatalf("provider merge is not idempotent: %s", contents)
	}
}

func TestQwenSettingsMergeRejectsInvalidJSON(t *testing.T) {
	settings := filepath.Join(t.TempDir(), "settings.json")
	invalid := "{broken\n"
	if err := os.WriteFile(settings, []byte(invalid), 0o600); err != nil {
		t.Fatal(err)
	}
	merger := filepath.Join("..", "ai", "install", "merge-qwen-settings.py")
	if output, err := exec.Command(merger, settings, "qwen3.5:4b", "16384").CombinedOutput(); err == nil {
		t.Fatalf("merger accepted invalid JSON: %s", output)
	}
	contents, err := os.ReadFile(settings)
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != invalid {
		t.Fatalf("merger changed invalid settings: %q", contents)
	}
	if _, err := os.Stat(settings + ".rapidou-backup"); !os.IsNotExist(err) {
		t.Fatalf("merger backed up or created data after rejecting invalid JSON")
	}
}

func writeExecutable(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(strings.TrimSpace(contents)+"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
}
