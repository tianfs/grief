package grief

import (
    "context"
    "github.com/gin-gonic/gin"
    v1 "github.com/tianfs/grief/internal/testdata/rpcpb/v1"
    tGrpc "github.com/tianfs/grief/transport/grpc"
    "github.com/tianfs/grief/transport/http"
    "google.golang.org/grpc"
    "log"
    "testing"
    "time"
)

type GreeterService struct {
    v1.UnimplementedGreeterServer
}

func (g *GreeterService) SayHello(ctx context.Context, r *v1.HelloRequest) (*v1.HelloReply, error) {
    return &v1.HelloReply{Message: " Server:" + r.Name}, nil
}
func TestApp(t *testing.T) {
    // http
    r := gin.Default()
    r.GET("/getToken", func(c *gin.Context) {
        c.JSON(999, 123123)
    })

    hs := http.NewServer(http.Router(r), http.Address(":8007"))
    // grpc
    gs := tGrpc.NewServer(tGrpc.Address(":9003"))
    v1.RegisterGreeterServer(gs, &GreeterService{})
    app := NewApp(
        ID("id123"),
        Name("grief"),
        Server(hs, gs),
    )
    go func() {
        time.Sleep(3 * time.Second)
        grpcClient()
    }()
    if err := app.Run(); err != nil {
        panic(err)
    }
    t.Logf("testapp %+v \n", app)

}

func grpcClient() {
    conn, err := grpc.Dial("127.0.0.1:9003", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("客户端连接异常: %v", err)
    }
    defer conn.Close()

    c := v1.NewGreeterClient(conn) // 获取GRPC句柄

    r, err := c.SayHello(context.Background(), &v1.HelloRequest{
        Name: "哈哈",
    })
    if err != nil {
        log.Fatalf("could not greet: %v", err)
    }
    log.Printf("####### get server Greeting response: %s", r.Message)

}
