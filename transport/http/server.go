package http

import (
    "context"
    "crypto/tls"
    "errors"
    "github.com/gin-gonic/gin"
    "github.com/tianfs/grief/transport"
    "log"
    "net"
    "net/http"
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

// Router is gin engine
func Router(ge *gin.Engine) ServerOption {
    return func(o *Server) {
        o.router = ge
    }
}

type Server struct {
    *http.Server
    network  string
    lis      net.Listener
    address  string
    endpoint *url.URL
    timeout  time.Duration
    router   *gin.Engine
    tlsConf  *tls.Config
}

func NewServer(opts ...ServerOption) *Server {
    srv := &Server{
        network: "tcp",
        address: ":8006",
        timeout: 1 * time.Second,
    }
    for _, opt := range opts {
        opt(srv)
    }

    srv.Server = &http.Server{
        Addr:    srv.address,
        Handler: srv.router,
    }
    srv.setListen()
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
    s.BaseContext = func(net.Listener) context.Context {
        return ctx
    }
    log.Printf("[HTTP] server listening on: %s \n", s.lis.Addr().String())
    var err error
    if s.tlsConf != nil {
        err = s.ServeTLS(s.lis, "", "")
    } else {
        err = s.Serve(s.lis)
    }
    if !errors.Is(err, http.ErrServerClosed) {
        return err
    }

    return nil
}

// Stop start the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
    log.Printf("执行http Shutdown\n")
    return s.Shutdown(ctx)
}
