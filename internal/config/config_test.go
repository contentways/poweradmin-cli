// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/contentways/poweradmin-cli/v2/internal/config"
)

func TestLoadMissingFile(t *testing.T) {
	cfg, err := config.Load("/tmp/nonexistent-poweradmin-config.yaml")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if cfg.URL != "" || cfg.APIKey != "" {
		t.Errorf("expected empty config, got: %+v", cfg)
	}
}

func TestLoadValidFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	content := "url: https://dns.example.com\napi_key: pwa_test\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.URL != "https://dns.example.com" {
		t.Errorf("URL = %q, want https://dns.example.com", cfg.URL)
	}
	if cfg.APIKey != "pwa_test" {
		t.Errorf("APIKey = %q, want pwa_test", cfg.APIKey)
	}
}

func TestDefaultPath(t *testing.T) {
	path := config.DefaultPath()
	if path == "" {
		t.Error("expected non-empty default path")
	}
}
