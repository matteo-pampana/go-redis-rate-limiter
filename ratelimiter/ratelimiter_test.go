package ratelimiter

import (
	"context"
	"errors"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter_CheckRequest(t *testing.T) {
	testCtx := context.Background()
	type fields struct {
		mockedStore func(ctrl *gomock.Controller) store
		config      RateLimiterConfig
	}
	type args struct {
		ctx  context.Context
		keys []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "happy path",
			fields: fields{
				mockedStore: func(ctrl *gomock.Controller) store {
					mockStore := NewMockstore(ctrl)
					mockStore.EXPECT().GetCounter(testCtx, "key1#key2").Return(0, nil)
					mockStore.EXPECT().IncreaseWithTTL(testCtx, "key1#key2", time.Second).Return(nil)
					return mockStore
				},
				config: RateLimiterConfig{
					MaxRequests:     10,
					RefreshInterval: time.Second,
				},
			},
			args: args{
				ctx:  testCtx,
				keys: []string{"key1", "key2"},
			},
			wantErr: false,
		},
		{
			name: "error when max requests reached",
			fields: fields{
				mockedStore: func(ctrl *gomock.Controller) store {
					mockStore := NewMockstore(ctrl)
					mockStore.EXPECT().GetCounter(testCtx, "key1#key2").Return(10, nil)
					return mockStore
				},
				config: RateLimiterConfig{
					MaxRequests:     10,
					RefreshInterval: time.Second,
				},
			},
			args: args{
				ctx:  testCtx,
				keys: []string{"key1", "key2"},
			},
			wantErr: true,
		},
		{
			name: "error when store returns error",
			fields: fields{
				mockedStore: func(ctrl *gomock.Controller) store {
					mockStore := NewMockstore(ctrl)
					mockStore.EXPECT().GetCounter(testCtx, "key1#key2").Return(0, errors.New("error"))
					return mockStore
				},
				config: RateLimiterConfig{
					MaxRequests:     10,
					RefreshInterval: time.Second,
				},
			},
			args: args{
				ctx:  testCtx,
				keys: []string{"key1", "key2"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockedStore := tt.fields.mockedStore(mockCtrl)
			r := NewRateLimiter(mockedStore, tt.fields.config)
			err := r.CheckRequest(tt.args.ctx, tt.args.keys)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestNewRateLimiter(t *testing.T) {
	type args struct {
		store  store
		config RateLimiterConfig
	}
	tests := []struct {
		name string
		args args
		want *RateLimiter
	}{
		{
			name: "happy path",
			args: args{
				store: NewMockstore(gomock.NewController(t)),
				config: RateLimiterConfig{
					MaxRequests:     10,
					RefreshInterval: time.Second,
				},
			},
			want: &RateLimiter{
				store: NewMockstore(gomock.NewController(t)),
				config: RateLimiterConfig{
					MaxRequests:     10,
					RefreshInterval: time.Second,
				},
			},
		},
		{
			name: "default config when not provided",
			args: args{
				store: NewMockstore(gomock.NewController(t)),
			},
			want: &RateLimiter{
				store: NewMockstore(gomock.NewController(t)),
				config: RateLimiterConfig{
					MaxRequests:     DEFAULT_MAX_REQUESTS,
					RefreshInterval: DEFAULT_REFRESH_INTERVAL,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			rl := NewRateLimiter(tt.args.store, tt.args.config)
			assert.Equal(t, tt.want, rl)
		})
	}
}
