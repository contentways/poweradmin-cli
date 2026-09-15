package records_test

import (
	"testing"

	"github.com/contentways/poweradmin-cli/v2/internal/cmd/records"
)

func TestNewRecordsCommand(t *testing.T) {
	cmd := records.NewRecordsCommand(nil)

	if cmd.Use != "records" {
		t.Fatalf("Use = %q", cmd.Use)
	}

	if len(cmd.Commands()) != 5 {
		t.Fatalf("got %d subcommands, want 5", len(cmd.Commands()))
	}
}
