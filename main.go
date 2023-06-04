package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/sys/unix"
	"html2pdf/convert"
	"net"
	"net/http"
	"os"
	"syscall"
)

//const replicaCount = 1
func main() {

	engine := gin.New()
	engine.UseH2C = true
	engine.SetTrustedProxies(nil)
	//engine.Use(gin.Logger(), gin.Recovery())
	engine.POST("/infra/pdfconv/internal/v1/convertHtml", wrap[convert.PdfRequest, *convert.ResultML[string]](convert.ConvertPdf))
	engine.GET("/test", func(c *gin.Context) {
		body := struct {
			Name string
			Pid  int
		}{
			Name: "haha",
			Pid:  os.Getpid(),
		}
		v := fmt.Sprintln(body)
		//context.JSON(http.StatusOK, &body)
		c.String(http.StatusOK, "%s", v)
	})

	var lc = net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			var opErr error
			if err := c.Control(func(fd uintptr) {
				opErr = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
			}); err != nil {
				return err
			}
			return opErr
		},
	}
	l, err := lc.Listen(context.Background(), "tcp", ":8080")
	if err != nil {
		panic(err)
	}

	//for i := 0; i < replicaCount; i++ {
	//	pid, _, errn := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	//	if errn != 0 {
	//		os.Exit(1)
	//	}
	//	//fmt.Println(eno)
	//	if pid == 0 {
	//		log.Println("i am child,ny pid:", os.Getpid())
	//		break
	//	} else {
	//		log.Printf("i am parent,my pid is:%d, my child id:%d\n", os.Getpid(), pid)
	//	}
	//}

	convert.InitPdf()
	engine.RunListener(l)

}

type WebHandler[T any, V any] func(req *T,pid int) V

func wrap[T any, V any](whFunc WebHandler[T, V]) gin.HandlerFunc {
	pid := os.Getpid()
	return func(context *gin.Context) {
		req := new(T)
		context.BindJSON(req)
		ret := whFunc(req,pid)
		context.JSON(200, ret)
	}
}
