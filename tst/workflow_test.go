package tst

import (
	"os"
	"os/exec"
	"path/filepath"
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

	hook := filepath.Join(root, "ai", "hooks", "prepare-agent.sh")
	cases := []struct {
		complexity string
		size       string
		model      string
		reasoning  string
	}{
		{"low", "small", "gpt-5.6-terra", "low"},
		{"low", "medium", "gpt-5.6-terra", "medium"},
		{"low", "large", "gpt-5.6-terra", "high"},
		{"medium", "small", "gpt-5.6-terra", "medium"},
		{"medium", "medium", "gpt-5.6-terra", "high"},
		{"medium", "large", "gpt-5.6-sol", "high"},
		{"high", "small", "gpt-5.6-sol", "xhigh"},
		{"high", "medium", "gpt-5.6-sol", "xhigh"},
		{"high", "large", "gpt-5.6-sol", "max"},
	}
	for _, test := range cases {
		output, err := exec.Command(hook, "select-agent-model", test.complexity, test.size).CombinedOutput()
		if err != nil {
			t.Fatalf("delegation hook failed for %s/%s: %v: %s", test.complexity, test.size, err, output)
		}
		expected := "model=" + test.model + "\nreasoning_effort=" + test.reasoning + "\n"
		if string(output) != expected {
			t.Fatalf("unexpected %s/%s selection: got %q, want %q", test.complexity, test.size, output, expected)
		}
	}

	if output, err := exec.Command(hook, "bypass", "low", "small").CombinedOutput(); err == nil {
		t.Fatalf("delegation hook accepted a spawn without the model skill: %s", output)
	}
	if output, err := exec.Command(hook, "select-agent-model", "unknown", "small").CombinedOutput(); err == nil {
		t.Fatalf("delegation hook accepted an unknown complexity: %s", output)
	}
	if output, err := exec.Command(hook, "select-agent-model", "low", "unknown").CombinedOutput(); err == nil {
		t.Fatalf("delegation hook accepted an unknown size: %s", output)
	}
}
