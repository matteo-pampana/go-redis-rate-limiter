package ratelimiter

import (
	context "context"
	"errors"
	"testing"
	"time"

	redis "github.com/go-redis/redis/v8"
	redismock "github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func TestStore_GetCounter(t *testing.T) {
	type args struct {
		ctx        context.Context
		requestKey string
	}
	tests := []struct {
		name      string
		mockRedis func() *redis.Client
		args      args
		want      int
		wantErr   bool
	}{
		{
			name: "happy path",
			mockRedis: func() *redis.Client {
				db, mock := redismock.NewClientMock()
				mock.ExpectGet("key1#key2").SetVal("0")
				return db
			},
			args: args{
				ctx:        context.Background(),
				requestKey: "key1#key2",
			},
			want: 0,
		},
		{
			name: "when redis return nil, return 0",
			mockRedis: func() *redis.Client {
				db, mock := redismock.NewClientMock()
				mock.ExpectGet("key1#key2").SetErr(redis.Nil)
				return db
			},
			args: args{
				ctx:        context.Background(),
				requestKey: "key1#key2",
			},
			want: 0,
		},
		{
			name: "error when redis get failed",
			mockRedis: func() *redis.Client {
				db, mock := redismock.NewClientMock()
				mock.ExpectGet("key1#key2").SetErr(errors.New("error"))
				return db
			},
			args: args{
				ctx:        context.Background(),
				requestKey: "key1#key2",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := NewStore(tt.mockRedis())
			val, err := s.GetCounter(tt.args.ctx, tt.args.requestKey)
			assert.Equal(t, tt.want, val)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestStore_IncreaseWithTTL(t *testing.T) {
	testCtx := context.Background()
	type args struct {
		ctx        context.Context
		requestKey string
		ttl        time.Duration
	}
	tests := []struct {
		name      string
		mockRedis func() (*redis.Client, redismock.ClientMock)
		args      args
		wantErr   bool
	}{
		{
			name: "happy path",
			mockRedis: func() (*redis.Client, redismock.ClientMock) {
				db, mock := redismock.NewClientMock()
				mock.ExpectIncr("key1#key2").SetVal(1)
				mock.ExpectExpire("key1#key2", time.Second).SetVal(true)
				return db, mock
			},
			args: args{
				ctx:        testCtx,
				requestKey: "key1#key2",
				ttl:        time.Second,
			},
		},
		{
			name: "error when redis incr failed",
			mockRedis: func() (*redis.Client, redismock.ClientMock) {
				db, mock := redismock.NewClientMock()
				mock.ExpectIncr("key1#key2").SetErr(errors.New("error"))
				return db, mock
			},
			args: args{
				ctx:        testCtx,
				requestKey: "key1#key2",
				ttl:        time.Second,
			},
			wantErr: true,
		},
		{
			name: "error when redis expire failed",
			mockRedis: func() (*redis.Client, redismock.ClientMock) {
				db, mock := redismock.NewClientMock()
				mock.ExpectIncr("key1#key2").SetVal(1)
				mock.ExpectExpire("key1#key2", time.Second).SetErr(errors.New("error"))
				return db, mock
			},
			args: args{
				ctx:        testCtx,
				requestKey: "key1#key2",
				ttl:        time.Second,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db, mock := tt.mockRedis()
			s := NewStore(db)
			err := s.IncreaseWithTTL(tt.args.ctx, tt.args.requestKey, tt.args.ttl)
			assert.Equal(t, tt.wantErr, err != nil)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
