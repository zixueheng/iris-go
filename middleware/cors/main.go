package main

import (
	"github.com/kataras/iris/v12"

	"github.com/iris-contrib/middleware/cors"
)

// api 跨域控制

func main() {
	app := iris.New()

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowCredentials: false,
	})

	v1 := app.Party("/api/v1", crs).AllowMethods(iris.MethodOptions) // <- important for the preflight.
	{
		v1.Get("/home", func(ctx iris.Context) {
			ctx.WriteString("Hello from /home")
		})
		v1.Get("/about", func(ctx iris.Context) {
			// ctx.WriteString("Hello from /about")
			ctx.JSON(iris.Map{"message": "ok"})
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

	app.Run(iris.Addr("localhost:8080"))
}
