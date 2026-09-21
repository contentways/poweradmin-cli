// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package version

import (
	"runtime/debug"
	"testing"
)

func TestDetectFromBuildInfoUsesMainVersion(t *testing.T) {
	info := &debug.BuildInfo{
		Main: debug.Module{Version: "v1.2.3"},
	}
	v, commit := detectFromBuildInfo(info)
	if v != "v1.2.3" {
		t.Errorf("expected version v1.2.3, got %q", v)
	}
	if commit != "" {
		t.Errorf("expected empty commit, got %q", commit)
	}
}

func TestDetectFromBuildInfoIgnoresDevelVersion(t *testing.T) {
	info := &debug.BuildInfo{
		Main: debug.Module{Version: "(devel)"},
	}
	v, _ := detectFromBuildInfo(info)
	if v != "" {
		t.Errorf("expected empty version for (devel), got %q", v)
	}
}

func TestDetectFromBuildInfoExtractsShortRevision(t *testing.T) {
	info := &debug.BuildInfo{
		Settings: []debug.BuildSetting{
			{Key: "vcs.revision", Value: "abcdef1234567890"},
		},
	}
	_, commit := detectFromBuildInfo(info)
	if commit != "abcdef1" {
		t.Errorf("expected commit truncated to 7 chars, got %q", commit)
	}
}

func TestDetectFromBuildInfoShortRevisionUnchanged(t *testing.T) {
	info := &debug.BuildInfo{
		Settings: []debug.BuildSetting{
			{Key: "vcs.revision", Value: "abc123"},
		},
	}
	_, commit := detectFromBuildInfo(info)
	if commit != "abc123" {
		t.Errorf("expected short revision unchanged, got %q", commit)
	}
}

func TestDetectFromBuildInfoNoVCSSettings(t *testing.T) {
	info := &debug.BuildInfo{}
	v, commit := detectFromBuildInfo(info)
	if v != "" || commit != "" {
		t.Errorf("expected both empty with no data, got v=%q commit=%q", v, commit)
	}
}

func TestDetectFromBuildInfoIgnoresOtherSettings(t *testing.T) {
	info := &debug.BuildInfo{
		Settings: []debug.BuildSetting{
			{Key: "GOOS", Value: "linux"},
			{Key: "vcs.revision", Value: "1234567"},
			{Key: "vcs.time", Value: "2026-01-01"},
		},
	}
	_, commit := detectFromBuildInfo(info)
	if commit != "1234567" {
		t.Errorf("expected commit 1234567, got %q", commit)
	}
}
