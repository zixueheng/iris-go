package main

import (
	"github.com/kataras/iris/v12"
)

// 跨域处理 普通的 发送 header 方式
// 或者使用 github.com/iris-contrib/middleware/cors 包：
// crs := cors.New(cors.Options{
// 	AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
// 	AllowCredentials: false,
// })
// v1 := app.Party("/api/v1", crs).AllowMethods(iris.MethodOptions)

func main() {
	app := iris.New()

	crs := func(ctx iris.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Origin,Content-Type")
		ctx.Next()
	} // or	"github.com/iris-contrib/middleware/cors"

	v1 := app.Party("/api/v1", crs).AllowMethods(iris.MethodOptions) // <- important for the preflight.
	{
		v1.Post("/mailer", func(ctx iris.Context) {
			var any iris.Map
			err := ctx.ReadJSON(&any)
			if err != nil {
				ctx.WriteString(err.Error())
				ctx.StatusCode(iris.StatusBadRequest)
				return
			}
			ctx.Application().Logger().Infof("received %#+v", any)
		})

		v1.Get("/home", func(ctx iris.Context) {
			ctx.WriteString("Hello from /home")
		})
		v1.Get("/about", func(ctx iris.Context) {
			ctx.WriteString("Hello from /about")
		})
		v1.Post("/send", func(ctx iris.Context) {
			ctx.WriteString("sent")
		})
		v1.Put("/send", func(ctx iris.Context) {
			ctx.WriteString("updated")
		})
		v1.Delete("/send", func(ctx iris.Context) {
			ctx.WriteString("deleted")
		})
	}

	// iris.WithoutPathCorrectionRedirection | iris#Configuration.DisablePathCorrectionRedirection:
	// CORS needs the allow origin headers in the redirect response as well, we have a solution for this:
	// Add the iris.WithoutPathCorrectionRedirection option
	// to directly fire the handler instead of redirection (which is the default behavior)
	// on request paths like "/v1/mailer/" when "/v1/mailer" route handler is registered.
	app.Listen(":80", iris.WithoutPathCorrectionRedirection)
}
