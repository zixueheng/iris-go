package main

import (
	"github.com/kataras/iris/v12"
)

// 下载文件，可用于大文件

func main() {
	app := iris.New()

	app.Get("/", func(ctx iris.Context) {
		file := "./1.zip"           // 源文件
		ctx.SendFile(file, "c.zip") // 客户端下载得到的文件名 c.zip
	})

	app.Listen(":8080")
}
