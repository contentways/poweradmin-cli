// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package schema_test

import (
	"testing"

	"github.com/contentways/poweradmin-cli/v3/internal/schema"
	"github.com/contentways/poweradmin-go/v4/poweradmin"
)

func TestDNSSECKeyFromSDKNilDS(t *testing.T) {
	got := schema.DNSSECKeyFromSDK(&poweradmin.DNSSECKey{ID: 1, Type: poweradmin.DNSSECKeyTypeZSK})

	if got.DS == nil || len(got.DS) != 0 {
		t.Errorf("DS = %v, want empty non-nil slice so JSON shows []", got.DS)
	}
	if got.Type != "zsk" {
		t.Errorf("Type = %q, want zsk", got.Type)
	}
}

func TestZoneDNSSECFromSDK(t *testing.T) {
	got := schema.ZoneDNSSECFromSDK("example.com", &poweradmin.ZoneDNSSEC{
		Enabled:   true,
		Presigned: true,
		DSRecords: []poweradmin.DSRecord{{KeyTag: 1, Algorithm: 13, DigestType: 2, Digest: "ab"}},
	})

	if got.Zone != "example.com" || !got.Enabled || !got.Presigned || len(got.DSRecords) != 1 || got.DSRecords[0].Digest != "ab" {
		t.Errorf("unexpected conversion: %+v", got)
	}
}

func TestServerStatusFromSDK(t *testing.T) {
	got := schema.ServerStatusFromSDK(&poweradmin.ServerStatus{
		Running: true,
		Version: "4.9.17",
		Slaves:  []poweradmin.SlaveStatus{{IP: "192.0.2.1", Status: "ok"}},
	})

	if !got.Running || got.Version != "4.9.17" || len(got.Slaves) != 1 || got.Slaves[0].IP != "192.0.2.1" {
		t.Errorf("unexpected conversion: %+v", got)
	}
}

func TestDNSSECKeyListFromSDK(t *testing.T) {
	got := schema.DNSSECKeyListFromSDK("example.com", []*poweradmin.DNSSECKey{
		{ID: 1, Type: poweradmin.DNSSECKeyTypeKSK, DS: []string{"1 13 2 ab"}},
		{ID: 2, Type: poweradmin.DNSSECKeyTypeZSK},
	})

	if got.Zone != "example.com" || got.Count != 2 || got.Keys[0].DS[0] != "1 13 2 ab" || got.Keys[1].Type != "zsk" {
		t.Errorf("unexpected conversion: %+v", got)
	}
}
