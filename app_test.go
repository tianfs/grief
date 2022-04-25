package grief

import (
    "github.com/gin-gonic/gin"
    "github.com/tianfs/grief/transport/http"
    "testing"
)

func TestApp(t *testing.T) {
    return
    r := gin.Default()
    r.GET("/getToken", func(c *gin.Context) {
        c.JSON(999, 123123)
    })

    hs := http.NewServer(http.Router(r), http.Address(":8007"))
    app := NewApp(
        ID("id123"),
        Name("grief"),
        Server(hs),
    )
    if err := app.Run(); err != nil {
        panic(err)
    }
    t.Logf("testapp %+v \n", app)

}
