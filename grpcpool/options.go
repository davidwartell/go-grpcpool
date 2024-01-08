package grpcpool

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"time"
)

type Options struct {
	keepalive             *keepalive.ClientParameters
	pingFunc              PingFunc
	withSnappyCompression bool
	withOtelTracing       bool
}

// PingFunc should send a GRPC ping/pong to the other side of conn.  Returns err or latency.
type PingFunc func(ctx context.Context, conn *grpc.ClientConn) (time.Duration, error)

type Option func(o *Options)

//goland:noinspection GoUnusedExportedFunction
func WithKeepaliveClientParams(keepalive keepalive.ClientParameters) Option {
	return func(o *Options) {
		o.keepalive = &keepalive
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithPingTestEveryNewConnection(pingFunc PingFunc) Option {
	return func(o *Options) {
		o.pingFunc = pingFunc
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithSnappyCompression(enabled bool) Option {
	return func(o *Options) {
		o.withSnappyCompression = enabled
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithOtelTracing(enabled bool) Option {
	return func(o *Options) {
		o.withOtelTracing = enabled
	}
}
