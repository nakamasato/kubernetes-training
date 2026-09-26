package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveVersions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sample.yaml")
	for _, tc := range []struct {
		name, content, pattern, want string
		fail                         bool
	}{
		{"pin", "image: test:v2.3.4", `test:(v[0-9.]+)`, "v2.3.4", false},
		{"same pin twice", "test:v2.3.4 test:v2.3.4", `test:(v[0-9.]+)`, "v2.3.4", false},
		{"conflicting pins", "test:v2.3.4 test:v3.0.0", `test:(v[0-9.]+)`, "", true},
		{"missing pin", "test:latest", `test:(v[0-9.]+)`, "", true},
		{"invalid regex", "test:v2.3.4", `(`, "", true},
		{"missing capture", "test:v2.3.4", `test:v[0-9.]+`, "", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := os.WriteFile(path, []byte(tc.content), 0600); err != nil {
				t.Fatal(err)
			}
			input := []Version{{Name: "sample", Source: &VersionSource{Path: path, Pattern: tc.pattern}}}
			got, err := resolveVersions(input)
			if (err != nil) != tc.fail {
				t.Fatalf("resolveVersions error = %v", err)
			}
			if !tc.fail && got[0].Version != tc.want {
				t.Fatalf("got %q, want %q", got[0].Version, tc.want)
			}
			if input[0].Version != "" {
				t.Fatal("mutated source configuration")
			}
		})
	}
}

func TestInvalidSources(t *testing.T) {
	for _, input := range [][]Version{
		{{Name: "missing", Source: &VersionSource{Path: "does-not-exist", Pattern: `(v[0-9.]+)`}}},
		{{Name: "duplicate"}, {Name: "duplicate"}},
		{{Name: "ambiguous", Version: "v1", Source: &VersionSource{Path: "irrelevant"}}},
	} {
		if _, err := resolveVersions(input); err == nil {
			t.Fatalf("accepted invalid input: %+v", input)
		}
	}
}

func TestVersionSectionPreservesContentsAndIsIdempotent(t *testing.T) {
	original := "# Training\n\n# Versions\n\n1. kustomize: old\n1. PostgresOperator: old\n\n# Contents\n\nKeep this tutorial and its historical version v1.0.0.\n"
	versions := []Version{{Name: "Kustomize", Version: "kustomize/v5.8.1", RepoUrl: "https://example.com"}, {Name: "Postgres Operator", Version: "v2.0.2", RepoUrl: "https://example.com"}}
	got, err := replaceVersionSection(original, versions)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(got, original[strings.Index(original, "# Contents\n"):]) {
		t.Fatal("changed tutorial content")
	}
	if strings.Contains(got, "old") || !strings.Contains(got, "1. Postgres Operator: [v2.0.2]") {
		t.Fatal("stale versions remain")
	}
	again, err := replaceVersionSection(got, versions)
	if err != nil || again != got {
		t.Fatal("generation is not idempotent")
	}
	if _, err := replaceVersionSection("# Missing headers", versions); err == nil {
		t.Fatal("accepted missing section boundaries")
	}
}

func TestReleaseLinks(t *testing.T) {
	for _, tc := range []struct {
		v    Version
		want string
	}{
		{Version{RepoUrl: "https://example.com"}, "[not pinned](https://example.com)"},
		{Version{RepoUrl: "https://example.com", Version: "latest"}, "[latest](https://example.com/releases)"},
		{Version{RepoUrl: "https://example.com", Version: "13.2.2", TagPrefix: "v"}, "[13.2.2](https://example.com/releases/tag/v13.2.2)"},
		{Version{Version: "18.6", ReleaseURL: "https://www.postgresql.org/docs/release/"}, "[18.6](https://www.postgresql.org/docs/release/)"},
	} {
		if got := tc.v.link(); got != tc.want {
			t.Errorf("got %q, want %q", got, tc.want)
		}
	}
}
