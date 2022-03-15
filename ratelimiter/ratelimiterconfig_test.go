package ratelimiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiterConfig_IsEmpty(t *testing.T) {
	type fields struct {
		KeyItems        []string
		MaxRequests     int
		RefreshInterval time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "empty config returns true",
			want: true,
		},
		{
			name: "config with keys returns false",
			fields: fields{
				KeyItems: []string{"key1"},
			},
			want: false,
		},
		{
			name: "config with max requests but without keys returns true",
			fields: fields{
				MaxRequests: 20,
			},
			want: true,
		},
		{
			name: "config with interval and max requests but without keys returns true",
			fields: fields{
				MaxRequests:     15,
				RefreshInterval: 2 * time.Minute,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			cfg := RateLimiterConfig{
				KeyItems:        tt.fields.KeyItems,
				MaxRequests:     tt.fields.MaxRequests,
				RefreshInterval: tt.fields.RefreshInterval,
			}
			isEmpty := cfg.IsEmpty()
			assert.Equal(t, tt.want, isEmpty)
		})
	}
}

func TestRateLimiterConfig_SetDefault(t *testing.T) {
	type fields struct {
		KeyItems        []string
		MaxRequests     int
		RefreshInterval time.Duration
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "happy path",
			fields: fields{
				KeyItems:        []string{"one", "two"},
				MaxRequests:     34,
				RefreshInterval: 1 * time.Second,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			cfg := RateLimiterConfig{
				KeyItems:        tt.fields.KeyItems,
				MaxRequests:     tt.fields.MaxRequests,
				RefreshInterval: tt.fields.RefreshInterval,
			}
			cfg.SetDefault()
			assert.Equal(t, 1, len(cfg.KeyItems))
			assert.Equal(t, DEFAULT_KEY, cfg.KeyItems[0])
			assert.Equal(t, DEFAULT_MAX_REQUESTS, cfg.MaxRequests)
			assert.Equal(t, DEFAULT_REFRESH_INTERVAL, cfg.RefreshInterval)
		})
	}
}
