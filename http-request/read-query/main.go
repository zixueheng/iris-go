// package main contains an example on how to use the ReadForm, but with the same way you can do the ReadJSON & ReadJSON
package main

import (
	"github.com/kataras/iris/v12"
)

// MyType ...
type MyType struct {
	Name string `url:"name"`
	Age  int    `url:"age"`
}

func main() {
	app := iris.New()

	app.Get("/", func(ctx iris.Context) {
		var t MyType
		err := ctx.ReadQuery(&t) // 读取url 参数 填充到 结构体
		if err != nil && !iris.IsErrPath(err) {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString(err.Error())
		}

		ctx.Writef("MyType: %#v", t)
	})

	// http://localhost:8080?name=iris&age=3
	app.Listen(":8080")
}
