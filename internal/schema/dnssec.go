// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package schema

import "github.com/contentways/poweradmin-go/v4/poweradmin"

// ZoneDNSSEC is the CLI output schema for the DNSSEC status of a zone.
type ZoneDNSSEC struct {
	Zone      string     `json:"zone" yaml:"zone"`
	Enabled   bool       `json:"enabled" yaml:"enabled"`
	Presigned bool       `json:"presigned" yaml:"presigned"`
	DSRecords []DSRecord `json:"ds_records" yaml:"ds_records"`
	DNSKey    *string    `json:"dnskey,omitempty" yaml:"dnskey,omitempty"`
}

// DSRecord is the CLI output schema for a DS record.
type DSRecord struct {
	KeyTag     int    `json:"keytag" yaml:"keytag"`
	Algorithm  int    `json:"algorithm" yaml:"algorithm"`
	DigestType int    `json:"digest_type" yaml:"digest_type"`
	Digest     string `json:"digest" yaml:"digest"`
}

// ZoneDNSSECFromSDK converts a poweradmin SDK ZoneDNSSEC to the CLI output schema.
func ZoneDNSSECFromSDK(zoneName string, d *poweradmin.ZoneDNSSEC) ZoneDNSSEC {
	records := make([]DSRecord, len(d.DSRecords))
	for i, ds := range d.DSRecords {
		records[i] = DSRecord{KeyTag: ds.KeyTag, Algorithm: ds.Algorithm, DigestType: ds.DigestType, Digest: ds.Digest}
	}
	return ZoneDNSSEC{
		Zone:      zoneName,
		Enabled:   d.Enabled,
		Presigned: d.Presigned,
		DSRecords: records,
		DNSKey:    d.DNSKey,
	}
}

// DNSSECKey is the CLI output schema for a DNSSEC key.
type DNSSECKey struct {
	ID          int      `json:"id" yaml:"id"`
	Type        string   `json:"type" yaml:"type"`
	KeyTag      int      `json:"keytag" yaml:"keytag"`
	Algorithm   string   `json:"algorithm,omitempty" yaml:"algorithm,omitempty"`
	AlgorithmID int      `json:"algorithm_id" yaml:"algorithm_id"`
	Bits        int      `json:"bits" yaml:"bits"`
	Active      bool     `json:"active" yaml:"active"`
	DNSKey      *string  `json:"dnskey,omitempty" yaml:"dnskey,omitempty"`
	DS          []string `json:"ds" yaml:"ds"`
}

// DNSSECKeyFromSDK converts a poweradmin SDK DNSSECKey to the CLI output schema.
func DNSSECKeyFromSDK(k *poweradmin.DNSSECKey) DNSSECKey {
	ds := k.DS
	if ds == nil {
		ds = []string{}
	}
	return DNSSECKey{
		ID:          k.ID,
		Type:        string(k.Type),
		KeyTag:      k.KeyTag,
		Algorithm:   k.Algorithm,
		AlgorithmID: k.AlgorithmID,
		Bits:        k.Bits,
		Active:      k.Active,
		DNSKey:      k.DNSKey,
		DS:          ds,
	}
}

// DNSSECKeyList wraps the DNSSEC keys of a zone in a root object for JSON output.
type DNSSECKeyList struct {
	Zone  string      `json:"zone" yaml:"zone"`
	Keys  []DNSSECKey `json:"keys" yaml:"keys"`
	Count int         `json:"count" yaml:"count"`
}

// DNSSECKeyListFromSDK converts a slice of SDK DNSSECKeys to the CLI output schema.
func DNSSECKeyListFromSDK(zoneName string, keys []*poweradmin.DNSSECKey) DNSSECKeyList {
	out := make([]DNSSECKey, len(keys))
	for i, k := range keys {
		out[i] = DNSSECKeyFromSDK(k)
	}
	return DNSSECKeyList{Zone: zoneName, Keys: out, Count: len(out)}
}
