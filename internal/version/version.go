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

	if v := info.Main.Version; v != "" && v != "(devel)" {
		Version = v
	}

	for _, s := range info.Settings {
		if s.Key == "vcs.revision" {
			rev := s.Value
			if len(rev) > 7 {
				rev = rev[:7]
			}
			Commit = rev
		}
	}
}
