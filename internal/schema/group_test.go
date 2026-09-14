package schema_test

import (
	"testing"

	"contentways.dev/contentways/poweradmin-go/v2/poweradmin"
	"github.com/contentways/poweradmin-cli/internal/schema"
)

func TestGroupFromSDK(t *testing.T) {
	g := &poweradmin.Group{
		ID:          1,
		Name:        "admins",
		Description: "Administrators",
		PermTemplID: 10,
		MemberCount: 5,
		ZoneCount:   20,
	}

	out := schema.GroupFromSDK(g)

	if out.ID != 1 ||
		out.Name != "admins" ||
		out.Description != "Administrators" ||
		out.PermTemplID != 10 ||
		out.MemberCount != 5 ||
		out.ZoneCount != 20 {
		t.Errorf("GroupFromSDK = %+v", out)
	}
}

func TestGroupListFromSDK(t *testing.T) {
	groups := []*poweradmin.Group{
		{
			ID:   1,
			Name: "admins",
		},
		{
			ID:   2,
			Name: "operators",
		},
	}

	out := schema.GroupListFromSDK(groups)

	if out.Count != 2 {
		t.Fatalf("Count = %d, want 2", out.Count)
	}

	if len(out.Groups) != 2 {
		t.Fatalf("Groups = %d, want 2", len(out.Groups))
	}
}
