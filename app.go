package grief

import (
    "context"
    "errors"
    "github.com/tianfs/grief/transport"
    "golang.org/x/sync/errgroup"
    "log"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"
)

// AppInfo is application context value.
type AppInfo interface {
    ID() string
    Name() string
}
type Option func(a *App)

// ID with service id.
func ID(id string) Option {
    return func(a *App) { a.id = id }
}

// Name with service name.
func Name(name string) Option {
    return func(a *App) { a.name = name }
}

// Server with transport servers.
func Server(srv ...transport.Server) Option {
    return func(a *App) { a.servers = srv }
}

// App is main
type App struct {
    id          string
    name        string
    rootCtx     context.Context
    ctx         context.Context
    sigs        []os.Signal
    stopTimeout time.Duration
    servers     []transport.Server
    cancel      func()
}

func NewApp(opts ...Option) *App {
    app := App{
        rootCtx:     context.Background(),
        sigs:        []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
        stopTimeout: 10 * time.Second,
    }
    for _, opt := range opts {
        opt(&app)
    }
    ctx, cancel := context.WithCancel(app.rootCtx)
    app.ctx = ctx
    app.cancel = cancel
    return &app

}

// ID returns app instance id.
func (a *App) ID() string { return a.id }

// Name returns service name.
func (a *App) Name() string { return a.name }

func (a *App) Run() error {
    log.Printf("`app RUn`122%+v \n", a)
    ctx := NewContext(a.ctx, a)
    eg, ctx := errgroup.WithContext(ctx)
    wg := sync.WaitGroup{}
    for _, srv := range a.servers {
        srv := srv
        eg.Go(func() error {
            <-ctx.Done()
            log.Printf("接收ctx.done到应用停止,开始退出srv\n")
            stopCtx, stopCancel := context.WithTimeout(NewContext(context.Background(), a), a.stopTimeout)
            defer stopCancel()
            return srv.Stop(stopCtx)

        })
        wg.Add(1)
        eg.Go(func() error {
            wg.Done()
            startCtx := NewContext(context.Background(), a)
            return srv.Start(startCtx)
        })
    }
    wg.Wait()
    // 服务之策
    // 服务监听推出信号
    ch := make(chan os.Signal, 1)
    signal.Notify(ch, a.sigs...)
    eg.Go(func() error {
        for {
            select {
            case <-ch:
                log.Printf("ch管道接收到推出信号，开始结束应用\n")
                return a.Stop()
            case <-ctx.Done():
                log.Printf("接收到上下问退出\n")
                return ctx.Err()
            }
        }
    })
    if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
        return err
    }
    return nil
}
func (a *App) Stop() error {
    log.Printf("执行应用停止\n")
    if a.cancel != nil {
        a.cancel()
    }
    return nil

}

type appKey struct{}

// NewContext returns a new Context that carries value.
func NewContext(ctx context.Context, s AppInfo) context.Context {
    return context.WithValue(ctx, appKey{}, s)
}

// FromContext returns the Transport value stored in ctx, if any.
func FromContext(ctx context.Context) (s AppInfo, ok bool) {
    s, ok = ctx.Value(appKey{}).(AppInfo)
    return
}
