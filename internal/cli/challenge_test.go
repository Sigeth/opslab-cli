package cli

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const validManifest = `apiVersion: opslab.io/v1alpha1
kind: ChallengePackage
metadata: {name: hello, version: 0.1.0}
spec:
  title: Hello
  summary: Smallest valid challenge.
  category: linux-hardening
  difficulty: easy
  machines:
    - {name: box, runtime: pod, image: docker.io/library/debian:13-slim, setup: [setup.sh]}
  access:
    terminals: [{machine: box}]
  tasks:
    - id: hello
      title: Create /tmp/hello
      points: 1
      check: {machine: box, file: {path: /tmp/hello}}
`

func run(t *testing.T, args ...string) (string, error) {
	t.Helper()
	cmd := NewRootCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return out.String(), err
}

func writePackage(t *testing.T, manifest string, files ...string) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "challenge.yaml"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(dir, f), nil, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestLintValid(t *testing.T) {
	out, err := run(t, "challenge", "lint", writePackage(t, validManifest, "setup.sh"))
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if !strings.Contains(out, "hello 0.1.0: valid (1 machines, 1 tasks)") {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestLintReportsEveryProblem(t *testing.T) {
	manifest := strings.Replace(validManifest, "difficulty: easy", "difficulty: trivial", 1)
	out, err := run(t, "challenge", "lint", writePackage(t, manifest)) // setup.sh missing too
	if !errors.Is(err, errInvalid) {
		t.Fatalf("expected errInvalid, got %v", err)
	}
	for _, want := range []string{"spec.difficulty", "spec.machines[0].setup[0]"} {
		if !strings.Contains(out, want) {
			t.Errorf("output does not mention %s:\n%s", want, out)
		}
	}
}

func TestLintMissingManifest(t *testing.T) {
	if _, err := run(t, "challenge", "lint", t.TempDir()); err == nil {
		t.Fatal("expected an error")
	}
}

func TestVersion(t *testing.T) {
	out, err := run(t, "version")
	if err != nil || strings.TrimSpace(out) == "" {
		t.Fatalf("version: %q, %v", out, err)
	}
}
