package cli

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"

	"github.com/Sigeth/opslab-core-go/api/v1alpha1"
)

func TestManifestPrintsConfigMapAndChallenge(t *testing.T) {
	dir := writePackage(t, validManifest, "setup.sh")
	if err := os.WriteFile(filepath.Join(dir, "setup.sh"), []byte("echo hi\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "instructions"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "instructions", "en.md"), []byte("# Hello\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	out, err := run(t, "challenge", "manifest", dir, "-n", "opslab-dev")
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	docs := strings.Split(out, "\n---\n")
	if len(docs) != 2 {
		t.Fatalf("expected 2 documents:\n%s", out)
	}
	var cm corev1.ConfigMap
	if err := yaml.UnmarshalStrict([]byte(docs[0]), &cm); err != nil {
		t.Fatal(err)
	}
	if cm.Namespace != "opslab-dev" || cm.Data["setup.sh"] != "echo hi\n" || cm.Data["instructions__en.md"] != "# Hello\n" {
		t.Fatalf("configmap: %+v", cm)
	}
	if _, ok := cm.Data["challenge.yaml"]; ok {
		t.Fatal("the manifest belongs in the Challenge, not in the ConfigMap")
	}
	var ch v1alpha1.Challenge
	if err := yaml.UnmarshalStrict([]byte(docs[1]), &ch); err != nil {
		t.Fatal(err)
	}
	if ch.Name != "hello" || ch.Spec.Source.Inline == nil || ch.Spec.Source.Inline.FilesConfigMap != cm.Name {
		t.Fatalf("challenge: %+v", ch)
	}
}

func TestManifestRefusesInvalidPackages(t *testing.T) {
	_, err := run(t, "challenge", "manifest", writePackage(t, validManifest)) // setup.sh missing
	if !errors.Is(err, errInvalid) {
		t.Fatalf("expected errInvalid, got %v", err)
	}
}

func TestImagesListsBuiltImages(t *testing.T) {
	dir := writePackage(t, validManifest, "setup.sh")
	out, err := run(t, "challenge", "images", dir)
	if err != nil || out != "box\tdocker.io/library/debian:13-slim\t-\n" {
		t.Fatalf("%v: %q", err, out)
	}
	if err := os.MkdirAll(filepath.Join(dir, "images", "box"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "images", "box", "Dockerfile"), []byte("FROM scratch\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, _ = run(t, "challenge", "images", dir)
	if want := "box\tdocker.io/library/debian:13-slim\t" + filepath.Join(dir, "images", "box") + "\n"; out != want {
		t.Fatalf("got %q, want %q", out, want)
	}
}
