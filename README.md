# go-grpcpool

GRPC client connection multiplexer

An interface for pooling GRPC client connections. With an implementation that multiplexes client requests to the server 
over a shared http2 socket.

* mutual-tls
* snappy compressor
* idle connection checking
* connecting testing and recovery
* connection watchdog
* open telemetry tracing across grpc

## Usage

Example
```go
func main() {
mutualTLSConnectionFactory, err := grpcpool.NewMutualTLSFactory(
[]byte("FIXME"),                        // PEM encoded CA cert string used to verify server cert signature during handshake
[]byte("FIXME"),                        // PEM encoded client certificate string
[]byte("FIXME"),                        // PEM encoded client private RSA key
"localhost:8080",                       // TCP connection string of server

// After a duration of this time if the pool doesn't see any activity it tests the connection with the server to 
// see if the transport is still alive. If set below 10s, a minimum value of 10s will be used instead.
time.Second*time.Duration(30),

// After having pinged for keepalive check, the client waits for a duration of this timeout and if no activity is seen 
// even after that the connection is forcibly closed.  Keep in mind on Linux & Windows you will never know if a 
// connection is alive until you try to write to it.
time.Second*time.Duration(10),

// pingConn function pointer to test GRPC connections (see below)
pingConn,
)
if err != nil {
err = errors.Errorf("error loading mutual TLS connection enrollmentFactory for site (enrollment): %v", err)
panic(err)
}

// you probably want one pool instance shared by your whole app
pool := shared.New(
// Your connectiuon factory.  I have supplied mutual-TLS, other connection types are possible by implementing: 
// type ConnectionFactory interface {
//      NewConnection(ctx context.Context) (*grpc.ClientConn, error)
//      ConnectionOk(ctx context.Context, conn *grpc.ClientConn) error
// }

// if GRPC connection is idle for longer than this it will be closed even if its good
time.Second*time.Duration(5*60)
)
// the pool should be closed on shutdown of your app
defer func() {
if pool != nil {
pool.Close()
pool = nil
}
}()

ctx := context.Background() // FIXME
var conn grpcpool.PoolClientConn
conn, err = pool.Get(ctx)
if err != nil {
panic(err)
}
defer func() {
if conn != nil {
conn.Return()
}
}()

// do something useful with the connection to send GRPC request with: conn.Connection()

return
}

// pingConn should implement a GRPC ping/pong to the other side of conn.  Returns err or latency.
func pingConn(ctx context.Context, conn *grpc.ClientConn) (time.Duration, error) {
panic("unimplemented")
}
```

## Contributing

Happy to accept PRs.

# Author

**davidwartell**

* <http://github.com/davidwartell>
* <http://linkedin.com/in/wartell>
