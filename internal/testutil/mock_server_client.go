// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package testutil

import (
	"context"

	"github.com/contentways/poweradmin-go/v4/poweradmin"
)

// MockServerClient implements poweradmin.IServerClient for testing.
type MockServerClient struct {
	StatusFn func(ctx context.Context, opts poweradmin.ServerStatusOpts) (*poweradmin.ServerStatus, *poweradmin.Response, error)
}

func (m *MockServerClient) Status(ctx context.Context, opts poweradmin.ServerStatusOpts) (*poweradmin.ServerStatus, *poweradmin.Response, error) {
	if m.StatusFn != nil {
		return m.StatusFn(ctx, opts)
	}
	return nil, nil, nil
}
