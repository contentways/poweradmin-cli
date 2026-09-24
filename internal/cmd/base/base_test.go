// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package base_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/cmd/base"
	"github.com/contentways/poweradmin-cli/v3/internal/output"
	"github.com/contentways/poweradmin-cli/v3/internal/testutil"
	"github.com/contentways/poweradmin-go/v4/poweradmin"
	"github.com/spf13/cobra"
)

// newTestCmd returns a bare *cobra.Command with its stdout wired to buf,
// suitable for passing to base package functions that read/write via the
// Cobra command rather than directly through os.Stdin/Stdout.
func newTestCmd(buf *strings.Builder) *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	cmd.SetOut(buf)
	return cmd
}

// --- PrintJSON / PrintYAML / PrintFormatted -------------------------------

func TestPrintJSON(t *testing.T) {
	var buf strings.Builder
	cmd := newTestCmd(&buf)

	err := base.PrintJSON(cmd, map[string]any{"name": "example.com", "id": 42})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), `"name": "example.com"`) {
		t.Errorf("expected JSON output to contain name, got:\n%s", buf.String())
	}
}

func TestPrintJSONMarshalError(t *testing.T) {
	var buf strings.Builder
	cmd := newTestCmd(&buf)

	// Channels cannot be marshaled to JSON.
	err := base.PrintJSON(cmd, map[string]any{"bad": make(chan int)})
	if err == nil {
		t.Fatal("expected error for unmarshalable value, got nil")
	}
}

func TestPrintYAML(t *testing.T) {
	var buf strings.Builder
	cmd := newTestCmd(&buf)

	err := base.PrintYAML(cmd, map[string]any{"name": "example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "name: example.com") {
		t.Errorf("expected YAML output to contain name, got:\n%s", buf.String())
	}
}

func TestPrintFormattedDispatchesToJSON(t *testing.T) {
	var buf strings.Builder
	cmd := newTestCmd(&buf)

	err := base.PrintFormatted(cmd, output.FormatJSON, map[string]any{"name": "example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), `"name": "example.com"`) {
		t.Errorf("expected JSON-shaped output, got:\n%s", buf.String())
	}
}

func TestPrintFormattedDispatchesToYAML(t *testing.T) {
	var buf strings.Builder
	cmd := newTestCmd(&buf)

	err := base.PrintFormatted(cmd, output.FormatYAML, map[string]any{"name": "example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "name: example.com") {
		t.Errorf("expected YAML-shaped output, got:\n%s", buf.String())
	}
}

// --- NewTable / IsQuiet ----------------------------------------------------

func TestNewTableWithHeader(t *testing.T) {
	var buf strings.Builder
	cmd := newTestCmd(&buf)
	cmd.Flags().Bool("no-header", false, "")

	tbl := base.NewTable(cmd)
	tbl.AddHeader("ID", "NAME")
	tbl.AddColoredRow(output.PlainCell("1"), output.PlainCell("example.com"))
	tbl.Flush()

	if !strings.Contains(buf.String(), "ID") {
		t.Errorf("expected header row in output, got:\n%s", buf.String())
	}
}

func TestNewTableNoHeader(t *testing.T) {
	var buf strings.Builder
	cmd := newTestCmd(&buf)
	cmd.Flags().Bool("no-header", false, "")
	if err := cmd.Flags().Set("no-header", "true"); err != nil {
		t.Fatalf("failed to set no-header flag: %v", err)
	}

	tbl := base.NewTable(cmd)
	tbl.AddHeader("ID", "NAME")
	tbl.AddColoredRow(output.PlainCell("1"), output.PlainCell("example.com"))
	tbl.Flush()

	if strings.Contains(buf.String(), "ID") {
		t.Errorf("expected header row to be suppressed, got:\n%s", buf.String())
	}
	if !strings.Contains(buf.String(), "example.com") {
		t.Errorf("expected data row to still be present, got:\n%s", buf.String())
	}
}

func TestIsQuietTrue(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().BoolP("quiet", "q", false, "")
	if err := cmd.Flags().Set("quiet", "true"); err != nil {
		t.Fatalf("failed to set quiet flag: %v", err)
	}
	if !base.IsQuiet(cmd) {
		t.Error("expected IsQuiet to return true")
	}
}

func TestIsQuietFalseByDefault(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().BoolP("quiet", "q", false, "")
	if base.IsQuiet(cmd) {
		t.Error("expected IsQuiet to return false by default")
	}
}

// --- Confirm ----------------------------------------------------------------

func TestConfirmYesFlagSkipsPrompt(t *testing.T) {
	var buf strings.Builder
	cmd := newTestCmd(&buf)
	cmd.Flags().BoolP("yes", "y", false, "")
	if err := cmd.Flags().Set("yes", "true"); err != nil {
		t.Fatalf("failed to set yes flag: %v", err)
	}

	if !base.Confirm(cmd, "Proceed? ") {
		t.Error("expected Confirm to return true when --yes is set")
	}
	if buf.String() != "" {
		t.Errorf("expected no prompt to be printed when --yes is set, got:\n%s", buf.String())
	}
}

func TestConfirmLowercaseY(t *testing.T) {
	var buf strings.Builder
	cmd := newTestCmd(&buf)
	cmd.Flags().BoolP("yes", "y", false, "")
	testutil.WithStdin(t, "y\n")

	if !base.Confirm(cmd, "Proceed? ") {
		t.Error("expected Confirm to return true for input 'y'")
	}
}

func TestConfirmUppercaseY(t *testing.T) {
	var buf strings.Builder
	cmd := newTestCmd(&buf)
	cmd.Flags().BoolP("yes", "y", false, "")
	testutil.WithStdin(t, "Y\n")

	if !base.Confirm(cmd, "Proceed? ") {
		t.Error("expected Confirm to return true for input 'Y'")
	}
}

func TestConfirmDeclined(t *testing.T) {
	var buf strings.Builder
	cmd := newTestCmd(&buf)
	cmd.Flags().BoolP("yes", "y", false, "")
	testutil.WithStdin(t, "n\n")

	if base.Confirm(cmd, "Proceed? ") {
		t.Error("expected Confirm to return false for input 'n'")
	}
	if !strings.Contains(buf.String(), "Aborted.") {
		t.Errorf("expected 'Aborted.' message, got:\n%s", buf.String())
	}
}

// --- ResolveZone -------------------------------------------------------------

func TestResolveZoneByName(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 1, Name: name}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("name", "", "")
	cmd.Flags().String("id", "", "")
	if err := cmd.Flags().Set("name", "example.com"); err != nil {
		t.Fatalf("failed to set name flag: %v", err)
	}
	cmd.SetContext(fx.State.WithContext(context.Background()))

	zone, err := base.ResolveZone(cmd, fx.State.MockClient)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if zone.Name != "example.com" {
		t.Errorf("expected zone name example.com, got %q", zone.Name)
	}
}

func TestResolveZoneByID(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByIDFn: func(ctx context.Context, id int) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: id, Name: "example.com"}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("name", "", "")
	cmd.Flags().String("id", "", "")
	if err := cmd.Flags().Set("id", "42"); err != nil {
		t.Fatalf("failed to set id flag: %v", err)
	}
	cmd.SetContext(fx.State.WithContext(context.Background()))

	zone, err := base.ResolveZone(cmd, fx.State.MockClient)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if zone.ID != 42 {
		t.Errorf("expected zone id 42, got %d", zone.ID)
	}
}

func TestResolveZoneInvalidID(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, &testutil.MockZoneClient{}, nil)

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("name", "", "")
	cmd.Flags().String("id", "", "")
	if err := cmd.Flags().Set("id", "not-a-number"); err != nil {
		t.Fatalf("failed to set id flag: %v", err)
	}
	cmd.SetContext(fx.State.WithContext(context.Background()))

	_, err := base.ResolveZone(cmd, fx.State.MockClient)
	if err == nil {
		t.Fatal("expected error for invalid id, got nil")
	}
}

func TestResolveZoneGetByNameError(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return nil, nil, fmt.Errorf("not found")
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("name", "", "")
	cmd.Flags().String("id", "", "")
	if err := cmd.Flags().Set("name", "missing.com"); err != nil {
		t.Fatalf("failed to set name flag: %v", err)
	}
	cmd.SetContext(fx.State.WithContext(context.Background()))

	_, err := base.ResolveZone(cmd, fx.State.MockClient)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestResolveZoneGetByIDError(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByIDFn: func(ctx context.Context, id int) (*poweradmin.Zone, *poweradmin.Response, error) {
			return nil, nil, fmt.Errorf("not found")
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("name", "", "")
	cmd.Flags().String("id", "", "")
	if err := cmd.Flags().Set("id", "1"); err != nil {
		t.Fatalf("failed to set id flag: %v", err)
	}
	cmd.SetContext(fx.State.WithContext(context.Background()))

	_, err := base.ResolveZone(cmd, fx.State.MockClient)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- ResolveUser -------------------------------------------------------------

func TestResolveUserByName(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		GetByNameFn: func(ctx context.Context, username string) (*poweradmin.User, *poweradmin.Response, error) {
			return &poweradmin.User{ID: 1, Username: username}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("name", "", "")
	cmd.Flags().String("id", "", "")
	if err := cmd.Flags().Set("name", "patrick"); err != nil {
		t.Fatalf("failed to set name flag: %v", err)
	}
	cmd.SetContext(fx.State.WithContext(context.Background()))

	user, err := base.ResolveUser(cmd, fx.State.MockClient)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Username != "patrick" {
		t.Errorf("expected username patrick, got %q", user.Username)
	}
}

func TestResolveUserByID(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		GetByIDFn: func(ctx context.Context, id int) (*poweradmin.User, *poweradmin.Response, error) {
			return &poweradmin.User{ID: id, Username: "patrick"}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("name", "", "")
	cmd.Flags().String("id", "", "")
	if err := cmd.Flags().Set("id", "3"); err != nil {
		t.Fatalf("failed to set id flag: %v", err)
	}
	cmd.SetContext(fx.State.WithContext(context.Background()))

	user, err := base.ResolveUser(cmd, fx.State.MockClient)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != 3 {
		t.Errorf("expected id 3, got %d", user.ID)
	}
}

func TestResolveUserInvalidID(t *testing.T) {
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, &testutil.MockUserClient{}, nil, nil)

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("name", "", "")
	cmd.Flags().String("id", "", "")
	if err := cmd.Flags().Set("id", "abc"); err != nil {
		t.Fatalf("failed to set id flag: %v", err)
	}
	cmd.SetContext(fx.State.WithContext(context.Background()))

	_, err := base.ResolveUser(cmd, fx.State.MockClient)
	if err == nil {
		t.Fatal("expected error for invalid id, got nil")
	}
}

func TestResolveUserGetByNameError(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		GetByNameFn: func(ctx context.Context, username string) (*poweradmin.User, *poweradmin.Response, error) {
			return nil, nil, fmt.Errorf("not found")
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("name", "", "")
	cmd.Flags().String("id", "", "")
	if err := cmd.Flags().Set("name", "ghost"); err != nil {
		t.Fatalf("failed to set name flag: %v", err)
	}
	cmd.SetContext(fx.State.WithContext(context.Background()))

	_, err := base.ResolveUser(cmd, fx.State.MockClient)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- ResolveGroup ------------------------------------------------------------

func TestResolveGroupByName(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Group, *poweradmin.Response, error) {
			return &poweradmin.Group{ID: 1, Name: name}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("name", "", "")
	cmd.Flags().String("id", "", "")
	if err := cmd.Flags().Set("name", "Editors"); err != nil {
		t.Fatalf("failed to set name flag: %v", err)
	}
	cmd.SetContext(fx.State.WithContext(context.Background()))

	group, err := base.ResolveGroup(cmd, fx.State.MockClient)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.Name != "Editors" {
		t.Errorf("expected name Editors, got %q", group.Name)
	}
}

func TestResolveGroupByID(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		GetByIDFn: func(ctx context.Context, id int) (*poweradmin.Group, *poweradmin.Response, error) {
			return &poweradmin.Group{ID: id, Name: "Editors"}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("name", "", "")
	cmd.Flags().String("id", "", "")
	if err := cmd.Flags().Set("id", "7"); err != nil {
		t.Fatalf("failed to set id flag: %v", err)
	}
	cmd.SetContext(fx.State.WithContext(context.Background()))

	group, err := base.ResolveGroup(cmd, fx.State.MockClient)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.ID != 7 {
		t.Errorf("expected id 7, got %d", group.ID)
	}
}

func TestResolveGroupInvalidID(t *testing.T) {
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, &testutil.MockGroupClient{}, nil)

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("name", "", "")
	cmd.Flags().String("id", "", "")
	if err := cmd.Flags().Set("id", "xyz"); err != nil {
		t.Fatalf("failed to set id flag: %v", err)
	}
	cmd.SetContext(fx.State.WithContext(context.Background()))

	_, err := base.ResolveGroup(cmd, fx.State.MockClient)
	if err == nil {
		t.Fatal("expected error for invalid id, got nil")
	}
}

func TestResolveGroupGetByNameError(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Group, *poweradmin.Response, error) {
			return nil, nil, fmt.Errorf("not found")
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("name", "", "")
	cmd.Flags().String("id", "", "")
	if err := cmd.Flags().Set("name", "Ghosts"); err != nil {
		t.Fatalf("failed to set name flag: %v", err)
	}
	cmd.SetContext(fx.State.WithContext(context.Background()))

	_, err := base.ResolveGroup(cmd, fx.State.MockClient)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- ResolveZoneID -----------------------------------------------------------

func TestResolveZoneIDFromZoneID(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, &testutil.MockZoneClient{}, nil)

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("zone-name", "", "")
	cmd.Flags().String("zone-id", "", "")
	if err := cmd.Flags().Set("zone-id", "99"); err != nil {
		t.Fatalf("failed to set zone-id flag: %v", err)
	}
	cmd.SetContext(fx.State.WithContext(context.Background()))

	id, err := base.ResolveZoneID(cmd, fx.State.MockClient)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 99 {
		t.Errorf("expected id 99, got %d", id)
	}
}

func TestResolveZoneIDFromZoneName(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return &poweradmin.Zone{ID: 55, Name: name}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("zone-name", "", "")
	cmd.Flags().String("zone-id", "", "")
	if err := cmd.Flags().Set("zone-name", "example.com"); err != nil {
		t.Fatalf("failed to set zone-name flag: %v", err)
	}
	cmd.SetContext(fx.State.WithContext(context.Background()))

	id, err := base.ResolveZoneID(cmd, fx.State.MockClient)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 55 {
		t.Errorf("expected id 55, got %d", id)
	}
}

func TestResolveZoneIDInvalidZoneID(t *testing.T) {
	fx := testutil.NewFixtureWithMocks(t, &testutil.MockZoneClient{}, nil)

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("zone-name", "", "")
	cmd.Flags().String("zone-id", "", "")
	if err := cmd.Flags().Set("zone-id", "not-a-number"); err != nil {
		t.Fatalf("failed to set zone-id flag: %v", err)
	}
	cmd.SetContext(fx.State.WithContext(context.Background()))

	_, err := base.ResolveZoneID(cmd, fx.State.MockClient)
	if err == nil {
		t.Fatal("expected error for invalid zone-id, got nil")
	}
}

func TestResolveZoneIDGetByNameError(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return nil, nil, fmt.Errorf("not found")
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("zone-name", "", "")
	cmd.Flags().String("zone-id", "", "")
	if err := cmd.Flags().Set("zone-name", "missing.com"); err != nil {
		t.Fatalf("failed to set zone-name flag: %v", err)
	}
	cmd.SetContext(fx.State.WithContext(context.Background()))

	_, err := base.ResolveZoneID(cmd, fx.State.MockClient)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- Completion functions ------------------------------------------------

func TestZoneNameCompletionFiltersByPrefix(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		AllFn: func(ctx context.Context) ([]*poweradmin.Zone, error) {
			return []*poweradmin.Zone{
				{Name: "example.com"},
				{Name: "example.org"},
				{Name: "other.net"},
			}, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	completion := base.ZoneNameCompletion(fx.State)
	names, directive := completion(&cobra.Command{Use: "test"}, nil, "example")

	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("unexpected directive: %v", directive)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 matches, got %d: %v", len(names), names)
	}
}

func TestZoneNameCompletionError(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		AllFn: func(ctx context.Context) ([]*poweradmin.Zone, error) {
			return nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	completion := base.ZoneNameCompletion(fx.State)
	names, directive := completion(&cobra.Command{Use: "test"}, nil, "")

	if directive != cobra.ShellCompDirectiveError {
		t.Errorf("expected error directive, got: %v", directive)
	}
	if names != nil {
		t.Errorf("expected nil names on error, got: %v", names)
	}
}

func TestUserNameCompletionFiltersByPrefix(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		AllFn: func(ctx context.Context) ([]*poweradmin.User, error) {
			return []*poweradmin.User{
				{Username: "patrick"},
				{Username: "paul"},
				{Username: "max"},
			}, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	completion := base.UserNameCompletion(fx.State)
	names, directive := completion(&cobra.Command{Use: "test"}, nil, "pa")

	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("unexpected directive: %v", directive)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 matches, got %d: %v", len(names), names)
	}
}

func TestGroupNameCompletionFiltersByPrefix(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		AllFn: func(ctx context.Context) ([]*poweradmin.Group, error) {
			return []*poweradmin.Group{
				{Name: "Editors"},
				{Name: "Everyone"},
				{Name: "Admins"},
			}, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)

	completion := base.GroupNameCompletion(fx.State)
	names, directive := completion(&cobra.Command{Use: "test"}, nil, "E")

	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("unexpected directive: %v", directive)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 matches, got %d: %v", len(names), names)
	}
}

func TestPermissionTemplateNameCompletionFiltersByPrefix(t *testing.T) {
	mockTmpl := &testutil.MockPermissionTemplateClient{
		ListFn: func(ctx context.Context) ([]*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return []*poweradmin.PermissionTemplate{
				{Name: "Zone Editors"},
				{Name: "Zone Admins"},
				{Name: "Read Only"},
			}, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, mockTmpl)

	completion := base.PermissionTemplateNameCompletion(fx.State)
	names, directive := completion(&cobra.Command{Use: "test"}, nil, "Zone")

	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("unexpected directive: %v", directive)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 matches, got %d: %v", len(names), names)
	}
}

func TestPermissionTemplateNameCompletionError(t *testing.T) {
	mockTmpl := &testutil.MockPermissionTemplateClient{
		ListFn: func(ctx context.Context) ([]*poweradmin.PermissionTemplate, *poweradmin.Response, error) {
			return nil, nil, fmt.Errorf("api error")
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, nil, mockTmpl)

	completion := base.PermissionTemplateNameCompletion(fx.State)
	names, directive := completion(&cobra.Command{Use: "test"}, nil, "")

	if directive != cobra.ShellCompDirectiveError {
		t.Errorf("expected error directive, got: %v", directive)
	}
	if names != nil {
		t.Errorf("expected nil names on error, got: %v", names)
	}
}

func TestResolveZoneNilWithoutError(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return nil, nil, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("name", "", "")
	cmd.Flags().String("id", "", "")
	if err := cmd.Flags().Set("name", "ghost.com"); err != nil {
		t.Fatalf("failed to set name flag: %v", err)
	}
	cmd.SetContext(fx.State.WithContext(context.Background()))

	_, err := base.ResolveZone(cmd, fx.State.MockClient)
	if err == nil {
		t.Fatal("expected error for nil zone with nil error, got nil")
	}
}

func TestResolveZoneIDNilWithoutError(t *testing.T) {
	mockZone := &testutil.MockZoneClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Zone, *poweradmin.Response, error) {
			return nil, nil, nil
		},
	}
	fx := testutil.NewFixtureWithMocks(t, mockZone, nil)

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("zone-name", "", "")
	cmd.Flags().String("zone-id", "", "")
	if err := cmd.Flags().Set("zone-name", "ghost.com"); err != nil {
		t.Fatalf("failed to set zone-name flag: %v", err)
	}
	cmd.SetContext(fx.State.WithContext(context.Background()))

	_, err := base.ResolveZoneID(cmd, fx.State.MockClient)
	if err == nil {
		t.Fatal("expected error for nil zone with nil error, got nil")
	}
}

func TestResolveUserNilWithoutError(t *testing.T) {
	mockUser := &testutil.MockUserClient{
		GetByNameFn: func(ctx context.Context, username string) (*poweradmin.User, *poweradmin.Response, error) {
			return nil, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, mockUser, nil, nil)

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("name", "", "")
	cmd.Flags().String("id", "", "")
	if err := cmd.Flags().Set("name", "ghost"); err != nil {
		t.Fatalf("failed to set name flag: %v", err)
	}
	cmd.SetContext(fx.State.WithContext(context.Background()))

	_, err := base.ResolveUser(cmd, fx.State.MockClient)
	if err == nil {
		t.Fatal("expected error for nil user with nil error, got nil")
	}
}

func TestResolveGroupNilWithoutError(t *testing.T) {
	mockGroup := &testutil.MockGroupClient{
		GetByNameFn: func(ctx context.Context, name string) (*poweradmin.Group, *poweradmin.Response, error) {
			return nil, nil, nil
		},
	}
	fx := testutil.NewFixtureWithAllMocks(t, nil, nil, nil, mockGroup, nil)

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("name", "", "")
	cmd.Flags().String("id", "", "")
	if err := cmd.Flags().Set("name", "Ghosts"); err != nil {
		t.Fatalf("failed to set name flag: %v", err)
	}
	cmd.SetContext(fx.State.WithContext(context.Background()))

	_, err := base.ResolveGroup(cmd, fx.State.MockClient)
	if err == nil {
		t.Fatal("expected error for nil group with nil error, got nil")
	}
}
