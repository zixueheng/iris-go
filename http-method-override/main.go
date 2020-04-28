package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/methodoverride"
	// "github.com/iris-contrib/middleware/cors"
)

// HTTP 方法覆盖是指 用 GET、POST等常用方法覆盖 PUT, DELETE 或者 PATCH
// 方法重写 包装器允许你在客户端不支持它的地方使用 HTTP 谓词，如 PUT 或 DELETE。
func main() {
	app := iris.New()

	// crs := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
	// 	AllowCredentials: false,
	// })
	// app.Use(crs)

	mo := methodoverride.New(
		// Defaults to nil.
		//
		methodoverride.SaveOriginalMethod("_originalMethod"),
		// Default values.
		//
		// methodoverride.Methods(http.MethodPost),
		// methodoverride.Headers("X-HTTP-Method",
		//                        "X-HTTP-Method-Override",
		//                        "X-Method-Override"),
		// methodoverride.FormField("_method"),
		// methodoverride.Query("_method"),
	)
	// Register it with `WrapRouter`.
	app.WrapRouter(mo)

	app.Post("/path", func(ctx iris.Context) {
		ctx.WriteString("post response")
	})

	app.Delete("/path", func(ctx iris.Context) {
		ctx.WriteString("delete response")
	})

	app.Run(iris.Addr(":8080"))
}
