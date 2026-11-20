// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package schema

import "github.com/contentways/poweradmin-go/v4/poweradmin"

// ServerStatus is the CLI output schema for the PowerDNS server status.
type ServerStatus struct {
	Running       bool              `json:"running" yaml:"running"`
	ServerID      string            `json:"server_id" yaml:"server_id"`
	DaemonType    string            `json:"daemon_type" yaml:"daemon_type"`
	Version       string            `json:"version" yaml:"version"`
	UptimeSeconds *int              `json:"uptime_seconds,omitempty" yaml:"uptime_seconds,omitempty"`
	Metrics       map[string]string `json:"metrics,omitempty" yaml:"metrics,omitempty"`
	Slaves        []SlaveStatus     `json:"slaves,omitempty" yaml:"slaves,omitempty"`
}

// SlaveStatus is the CLI output schema for an autoprimary server probe.
type SlaveStatus struct {
	IP          string  `json:"ip" yaml:"ip"`
	Status      string  `json:"status" yaml:"status"`
	LastChecked *string `json:"last_checked,omitempty" yaml:"last_checked,omitempty"`
	Error       *string `json:"error,omitempty" yaml:"error,omitempty"`
}

// ServerStatusFromSDK converts a poweradmin SDK ServerStatus to the CLI output schema.
func ServerStatusFromSDK(s *poweradmin.ServerStatus) ServerStatus {
	slaves := make([]SlaveStatus, len(s.Slaves))
	for i, sl := range s.Slaves {
		slaves[i] = SlaveStatus{IP: sl.IP, Status: sl.Status, LastChecked: sl.LastChecked, Error: sl.Error}
	}
	return ServerStatus{
		Running:       s.Running,
		ServerID:      s.ServerID,
		DaemonType:    s.DaemonType,
		Version:       s.Version,
		UptimeSeconds: s.UptimeSeconds,
		Metrics:       s.Metrics,
		Slaves:        slaves,
	}
}
