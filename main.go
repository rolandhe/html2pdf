package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html2pdf/convert"
)

func main() {

	var data []string

	//js := `["a","b"]`
	//json.Unmarshal([]byte(js), &data)
	data = []string{}
	if data == nil {
		fmt.Println("nil")
	}

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	engine.POST("/infra/pdfconv/internal/v1/convertHtml", wrap[convert.PdfRequest, *convert.ResultML[string]](convert.ConvertPdf))

	convert.InitPdf()
	engine.Run(":8080")
}

type WebHandler[T any, V any] func(req *T) V

func wrap[T any, V any](whFunc WebHandler[T, V]) gin.HandlerFunc {
	return func(context *gin.Context) {
		req := new(T)
		context.BindJSON(req)
		ret := whFunc(req)
		context.JSON(200, ret)
	}
}
