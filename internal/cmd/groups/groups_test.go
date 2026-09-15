package groups_test

import (
	"testing"

	"github.com/contentways/poweradmin-cli/v2/internal/cmd/groups"
)

func TestNewGroupsCommand(t *testing.T) {
	cmd := groups.NewGroupsCommand(nil)

	if cmd.Use != "groups" {
		t.Fatalf("Use = %q", cmd.Use)
	}

	if len(cmd.Commands()) != 11 {
		t.Fatalf("got %d subcommands, want 11", len(cmd.Commands()))
	}
}
