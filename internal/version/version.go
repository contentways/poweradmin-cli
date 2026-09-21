// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

// Package version holds build-time version information injected via ldflags.
package version

import "runtime/debug"

// Version and Commit are set at build time via GoReleaser ldflags.
// Defaults to "dev" and "none" for local builds; if unset, they are
// derived at runtime from Go's embedded module/VCS build info (e.g.
// when installed via `go install ...@version`).
var (
	Version = "dev"
	Commit  = "none"
)

func init() {
	if Version != "dev" || Commit != "none" {
		return
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	v, c := detectFromBuildInfo(info)
	if v != "" {
		Version = v
	}
	if c != "" {
		Commit = c
	}
}

// detectFromBuildInfo extracts a version string and a short (7-character)
// VCS revision from Go's embedded build info. Returns empty strings for
// either value that could not be determined, leaving the caller's existing
// default in place. Separated from init() so the detection logic itself can
// be tested with a fabricated *debug.BuildInfo, since init() cannot be
// re-triggered in tests.
func detectFromBuildInfo(info *debug.BuildInfo) (v, commit string) {
	if info.Main.Version != "" && info.Main.Version != "(devel)" {
		v = info.Main.Version
	}

	for _, s := range info.Settings {
		if s.Key == "vcs.revision" {
			rev := s.Value
			if len(rev) > 7 {
				rev = rev[:7]
			}
			commit = rev
		}
	}

	return v, commit
}
