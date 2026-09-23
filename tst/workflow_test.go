package tst

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestDelegationWorkflow(t *testing.T) {
	root := ".."
	for _, name := range []string{
		"select-agent-model", "luna", "spec", "questions", "backend", "frontend",
		"run-functional-tests", "code-review", "craft",
	} {
		path := filepath.Join(root, "ai", "skills", name, "SKILL.md")
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("required workflow skill %s is unavailable: %v", name, err)
		}
		for _, clientDirectory := range []string{".agents", ".claude"} {
			discovered := filepath.Join(root, clientDirectory, "skills", name, "SKILL.md")
			if _, err := os.Stat(discovered); err != nil {
				t.Fatalf("%s cannot discover skill %s: %v", clientDirectory, name, err)
			}
		}
	}
	for _, path := range []string{
		filepath.Join(root, ".codex", "hooks.json"),
		filepath.Join(root, ".claude", "settings.json"),
		filepath.Join(root, "CLAUDE.md"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("missing client integration %s: %v", path, err)
		}
	}

	policy, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{"Before spawning any independent agent", "select-agent-model/SKILL.md", "prepare-agent.sh"} {
		if !strings.Contains(string(policy), required) {
			t.Fatalf("AGENTS.md is missing mandatory delegation policy %q", required)
		}
	}

	selector := filepath.Join(root, "ai", "hooks", "prepare-agent.sh")
	selectionDir := t.TempDir()
	cases := []struct {
		platform   string
		complexity string
		size       string
		model      string
		reasoning  string
	}{
		{"codex", "low", "small", "gpt-5.6-terra", "low"},
		{"codex", "low", "medium", "gpt-5.6-terra", "medium"},
		{"codex", "low", "large", "gpt-5.6-terra", "high"},
		{"codex", "medium", "small", "gpt-5.6-terra", "medium"},
		{"codex", "medium", "medium", "gpt-5.6-terra", "high"},
		{"codex", "medium", "large", "gpt-5.6-sol", "high"},
		{"codex", "high", "small", "gpt-5.6-sol", "xhigh"},
		{"codex", "high", "medium", "gpt-5.6-sol", "xhigh"},
		{"codex", "high", "large", "gpt-5.6-sol", "max"},
		{"claude", "low", "small", "haiku", ""},
		{"claude", "low", "medium", "haiku", ""},
		{"claude", "low", "large", "sonnet", ""},
		{"claude", "medium", "small", "sonnet", ""},
		{"claude", "medium", "medium", "sonnet", ""},
		{"claude", "medium", "large", "opus", ""},
		{"claude", "high", "small", "opus", ""},
		{"claude", "high", "medium", "opus", ""},
		{"claude", "high", "large", "opus", ""},
	}
	for _, test := range cases {
		command := exec.Command(selector, "select-agent-model", test.platform, test.complexity, test.size)
		command.Env = append(os.Environ(), "RAPIDOU_SELECTION_DIR="+selectionDir)
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("delegation hook failed for %s/%s/%s: %v: %s", test.platform, test.complexity, test.size, err, output)
		}
		expected := "platform=" + test.platform + "\nmodel=" + test.model + "\n"
		if test.reasoning != "" {
			expected += "reasoning_effort=" + test.reasoning + "\n"
		}
		if string(output) != expected {
			t.Fatalf("unexpected %s/%s/%s selection: got %q, want %q", test.platform, test.complexity, test.size, output, expected)
		}
	}

	if output, err := exec.Command(selector, "bypass", "codex", "low", "small").CombinedOutput(); err == nil {
		t.Fatalf("delegation hook accepted a spawn without the model skill: %s", output)
	}
	if output, err := exec.Command(selector, "select-agent-model", "unknown", "low", "small").CombinedOutput(); err == nil {
		t.Fatalf("delegation hook accepted an unknown platform: %s", output)
	}
	if output, err := exec.Command(selector, "select-agent-model", "codex", "unknown", "small").CombinedOutput(); err == nil {
		t.Fatalf("delegation hook accepted an unknown complexity: %s", output)
	}
	if output, err := exec.Command(selector, "select-agent-model", "claude", "low", "unknown").CombinedOutput(); err == nil {
		t.Fatalf("delegation hook accepted an unknown size: %s", output)
	}

	testSpawnGuard(t, root, selector, selectionDir, "codex", "gpt-5.6-sol", "high", "medium", "xhigh")
	testSpawnGuard(t, root, selector, selectionDir, "claude", "sonnet", "medium", "medium", "")
}

func testSpawnGuard(t *testing.T, root, selector, selectionDir, platform, model, complexity, size, reasoning string) {
	t.Helper()
	selectCommand := exec.Command(selector, "select-agent-model", platform, complexity, size)
	selectCommand.Env = append(os.Environ(), "RAPIDOU_SELECTION_DIR="+selectionDir)
	if output, err := selectCommand.CombinedOutput(); err != nil {
		t.Fatalf("prepare %s selection: %v: %s", platform, err, output)
	}

	enforcer := filepath.Join(root, "ai", "hooks", "enforce-agent-selection.sh")
	wrong := exec.Command(enforcer, platform)
	wrong.Env = append(os.Environ(), "RAPIDOU_SELECTION_DIR="+selectionDir)
	wrong.Stdin = strings.NewReader(`{"tool_name":"spawn_agent","tool_input":{"model":"wrong"}}`)
	if output, err := wrong.CombinedOutput(); err == nil {
		t.Fatalf("%s guard accepted the wrong model: %s", platform, output)
	}
	selectCommand = exec.Command(selector, "select-agent-model", platform, complexity, size)
	selectCommand.Env = append(os.Environ(), "RAPIDOU_SELECTION_DIR="+selectionDir)
	if output, err := selectCommand.CombinedOutput(); err != nil {
		t.Fatalf("prepare replacement %s selection: %v: %s", platform, err, output)
	}

	payload := `{"tool_name":"spawn_agent","tool_input":{"model":"` + model + `"`
	if reasoning != "" {
		payload += `,"reasoning_effort":"` + reasoning + `"`
	}
	payload += `}}`
	allowed := exec.Command(enforcer, platform)
	allowed.Env = append(os.Environ(), "RAPIDOU_SELECTION_DIR="+selectionDir)
	allowed.Stdin = strings.NewReader(payload)
	if output, err := allowed.CombinedOutput(); err != nil {
		t.Fatalf("%s guard rejected the selected model: %v: %s", platform, err, output)
	}

	reused := exec.Command(enforcer, platform)
	reused.Env = append(os.Environ(), "RAPIDOU_SELECTION_DIR="+selectionDir)
	reused.Stdin = strings.NewReader(payload)
	if output, err := reused.CombinedOutput(); err == nil {
		t.Fatalf("%s guard reused a consumed selection: %s", platform, output)
	}

	staleState := "platform=" + platform + "\nmodel=" + model + "\nselected_at=1\n"
	if reasoning != "" {
		staleState += "reasoning_effort=" + reasoning + "\n"
	}
	if err := os.WriteFile(filepath.Join(selectionDir, "rapidou-agent-selection-"+platform), []byte(staleState), 0o600); err != nil {
		t.Fatal(err)
	}
	stale := exec.Command(enforcer, platform)
	stale.Env = append(os.Environ(), "RAPIDOU_SELECTION_DIR="+selectionDir)
	stale.Stdin = strings.NewReader(payload)
	if output, err := stale.CombinedOutput(); err == nil {
		t.Fatalf("%s guard accepted a stale selection: %s", platform, output)
	}

	selectCommand = exec.Command(selector, "select-agent-model", platform, complexity, size)
	selectCommand.Env = append(os.Environ(), "RAPIDOU_SELECTION_DIR="+selectionDir)
	if output, err := selectCommand.CombinedOutput(); err != nil {
		t.Fatalf("prepare concurrent %s selection: %v: %s", platform, err, output)
	}
	results := make(chan bool, 2)
	for range 2 {
		go func() {
			concurrent := exec.Command(enforcer, platform)
			concurrent.Env = append(os.Environ(), "RAPIDOU_SELECTION_DIR="+selectionDir)
			concurrent.Stdin = strings.NewReader(payload)
			_, err := concurrent.CombinedOutput()
			results <- err == nil
		}()
	}
	allowedCount := 0
	for range 2 {
		if <-results {
			allowedCount++
		}
	}
	if allowedCount != 1 {
		t.Fatalf("%s guard allowed %d concurrent spawns from one selection", platform, allowedCount)
	}
}

func TestInstall(t *testing.T) {
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	project := t.TempDir()
	if err := os.MkdirAll(filepath.Join(project, "lib"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(root, filepath.Join(project, "lib", "rapidou")); err != nil {
		t.Fatal(err)
	}

	installer := filepath.Join(project, "lib", "rapidou", "run", "install")
	for attempt := 0; attempt < 2; attempt++ {
		if output, err := exec.Command(installer, project).CombinedOutput(); err != nil {
			t.Fatalf("install attempt %d failed: %v: %s", attempt+1, err, output)
		}
	}

	for _, directory := range []string{".agents", ".claude"} {
		target, err := os.Readlink(filepath.Join(project, directory, "skills", "craft"))
		if err != nil {
			t.Fatalf("read %s craft link: %v", directory, err)
		}
		if target != "../../lib/rapidou/ai/skills/craft" {
			t.Fatalf("unexpected %s craft link: %s", directory, target)
		}
	}
	for path, target := range map[string]string{
		filepath.Join(".codex", "hooks.json"):     "../lib/rapidou/ai/install/codex-hooks.json",
		filepath.Join(".claude", "settings.json"): "../lib/rapidou/ai/install/claude-settings.json",
	} {
		got, err := os.Readlink(filepath.Join(project, path))
		if err != nil {
			t.Fatalf("read installed config %s: %v", path, err)
		}
		if got != target {
			t.Fatalf("unexpected config link %s: got %s, want %s", path, got, target)
		}
	}
	assertFileContains(t, filepath.Join(project, "AGENTS.md"), "lib/rapidou/docs/index.md")
	assertFileContains(t, filepath.Join(project, "CLAUDE.md"), "@AGENTS.md")

	conflictProject := t.TempDir()
	if err := os.MkdirAll(filepath.Join(conflictProject, "lib"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(root, filepath.Join(conflictProject, "lib", "rapidou")); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(conflictProject, ".codex"), 0o755); err != nil {
		t.Fatal(err)
	}
	conflict := filepath.Join(conflictProject, ".codex", "hooks.json")
	conflictContents := "{\"note\":\"lib/rapidou/ai/hooks/enforce-agent-selection.sh\"}\n"
	if err := os.WriteFile(conflict, []byte(conflictContents), 0o644); err != nil {
		t.Fatal(err)
	}
	if output, err := exec.Command(filepath.Join(conflictProject, "lib", "rapidou", "run", "install"), conflictProject).CombinedOutput(); err == nil {
		t.Fatalf("installer overwrote an unrelated config: %s", output)
	}
	contents, err := os.ReadFile(conflict)
	if err != nil || string(contents) != conflictContents {
		t.Fatalf("installer changed conflicting config: %q, %v", contents, err)
	}

	mergedProject := t.TempDir()
	for _, directory := range []string{filepath.Join("lib"), filepath.Join(".codex"), filepath.Join(".claude")} {
		if err := os.MkdirAll(filepath.Join(mergedProject, directory), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Symlink(root, filepath.Join(mergedProject, "lib", "rapidou")); err != nil {
		t.Fatal(err)
	}
	for source, destination := range map[string]string{
		filepath.Join(root, "ai", "install", "codex-hooks.json"):     filepath.Join(mergedProject, ".codex", "hooks.json"),
		filepath.Join(root, "ai", "install", "claude-settings.json"): filepath.Join(mergedProject, ".claude", "settings.json"),
	} {
		contents, err := os.ReadFile(source)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(destination, contents, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	mergedInstaller := filepath.Join(mergedProject, "lib", "rapidou", "run", "install")
	if output, err := exec.Command(mergedInstaller, "--accept-existing-hooks", mergedProject).CombinedOutput(); err != nil {
		t.Fatalf("installer rejected confirmed merged hooks: %v: %s", err, output)
	}
}

func TestDocumentationRoutes(t *testing.T) {
	root := ".."
	readme, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if len(readme) > 4000 {
		t.Fatalf("README is no longer a concise onboarding page: %d bytes", len(readme))
	}
	for _, required := range []string{
		"Any Application", "git submodule add", "./lib/rapidou/run/install",
		"craft skill", "docs/shared/principles.md", "AGENTS.md", "CLAUDE.md",
	} {
		if !strings.Contains(string(readme), required) {
			t.Fatalf("README is missing onboarding route %q", required)
		}
	}

	for path, required := range map[string][]string{
		"AGENTS.md":                   {"docs/index.md", "docs/shared/principles.md", "README.md"},
		"docs/shared/principles.md":   {"Direct architecture", "Frontend", "Functional testing", "simplest design"},
		"docs/example.md":             {"Museo Pixel", "admin@rapidou.test", "APP_JWT_SECRET"},
		"docs/specs/documentation.md": {"Happy path", "Expected mistakes", "Acceptance"},
	} {
		contents, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		for _, text := range required {
			if !strings.Contains(string(contents), text) {
				t.Fatalf("%s is missing %q", path, text)
			}
		}
	}

	for _, path := range []string{
		"README.md", "docs/index.md", "docs/installation.md", "docs/example.md",
		"docs/shared/index.md", "docs/shared/principles.md", "docs/specs/index.md",
		"docs/harness/index.md",
	} {
		assertLocalMarkdownLinksResolve(t, root, path)
	}
}

var markdownLink = regexp.MustCompile(`\[[^]]+\]\(([^)]+)\)`)

func assertLocalMarkdownLinksResolve(t *testing.T, root, path string) {
	t.Helper()
	contents, err := os.ReadFile(filepath.Join(root, path))
	if err != nil {
		t.Fatal(err)
	}
	for _, match := range markdownLink.FindAllStringSubmatch(string(contents), -1) {
		target := strings.Split(match[1], "#")[0]
		if target == "" || strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
			continue
		}
		resolved := filepath.Join(root, filepath.Dir(path), filepath.FromSlash(target))
		if _, err := os.Stat(resolved); err != nil {
			t.Fatalf("broken local link in %s: %s (%v)", path, match[1], err)
		}
	}
}

func assertFileContains(t *testing.T, path, expected string) {
	t.Helper()
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(contents), expected) {
		t.Fatalf("%s does not contain %q", path, expected)
	}
}
