package main

import "github.com/kataras/iris/v12"

// 响应记录器是一个 Iris 的特定 http.ResponseWriter ，它记录发送的正文， 状态代码和标头，你可以在路由的请求处理程序链中的任何处理程序上进行操作。

func main() {
	app := iris.New()

	app.Use(func(ctx iris.Context) {
		// 在发送数据之前调用  ctx.Record()
		ctx.Record() // 开始记录
		ctx.Next()   // 执行下一个
	})

	// 收集并“记录”。
	app.Done(func(ctx iris.Context) {
		// ctx.Recorder() 返回一个 ResponseRecorder. 其方法可用于操纵或检索响应。
		body := ctx.Recorder().Body()

		// 应该打印成功。
		app.Logger().Infof("sent: %s", string(body))
	})

	app.Get("/save", func(ctx iris.Context) {
		ctx.WriteString("success")
		ctx.Next() // calls the Done middleware(s).
	})

	app.Run(iris.Addr(":8080"))
}
