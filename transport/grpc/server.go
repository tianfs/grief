package grpc

import (
    "context"
    "crypto/tls"
    "github.com/tianfs/grief/transport"
    "google.golang.org/grpc"
    "google.golang.org/grpc/health"
    "google.golang.org/grpc/health/grpc_health_v1"
    "google.golang.org/grpc/reflection"
    "log"
    "net"
    "net/url"
    "time"
)

var _ transport.Server = (*Server)(nil)

type ServerOption func(*Server)

// Network with server network.
func Network(network string) ServerOption {
    return func(s *Server) {
        s.network = network
    }
}

// Address with server address.
func Address(address string) ServerOption {
    return func(s *Server) {
        s.address = address
    }
}

// Timeout with server timeout.
func Timeout(timeout time.Duration) ServerOption {
    return func(s *Server) {
        s.timeout = timeout
    }
}

// TLSConfig with TLS config.
func TLSConfig(c *tls.Config) ServerOption {
    return func(o *Server) {
        o.tlsConf = c
    }
}

type Server struct {
    *grpc.Server
    baseCtx  context.Context
    network  string
    lis      net.Listener
    address  string
    endpoint *url.URL
    timeout  time.Duration
    tlsConf  *tls.Config
    health   *health.Server
}

func NewServer(opts ...ServerOption) *Server {
    srv := &Server{
        network: "tcp",
        address: ":0",
        timeout: 1 * time.Second,
        health:  health.NewServer(),
    }
    for _, opt := range opts {
        opt(srv)
    }
    srv.Server = grpc.NewServer()

    srv.setListen()
    // 增加心跳检测
    grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
    reflection.Register(srv.Server)
    return srv
}
func (s *Server) setListen() error {
    if s.lis == nil {
        lis, err := net.Listen(s.network, s.address)
        if err != nil {
            return err
        }
        s.lis = lis
    }
    return nil
}

// Start start the HTTP server.
func (s *Server) Start(ctx context.Context) error {

    s.baseCtx = ctx
    log.Printf("[gRPC] server listening on: %s \n", s.lis.Addr().String())
    s.health.Resume()
    return s.Serve(s.lis)
}

// Stop start the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
    s.health.Shutdown()
    s.GracefulStop()
    log.Printf("[gRPC] server stopping")
    return nil
}
