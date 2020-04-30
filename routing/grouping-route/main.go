package main

import "github.com/kataras/iris/v12"

// 分组路由

// 分组路由可以使用相同的路由前缀，和同样的中间件和一样的布局模板

func main() {
	app := iris.Default()

	// users := app.Party("/users", myAuthMiddlewareHandler)
	// users.Get("/{id:uint64}/profile", userProfileHandler)  // http://localhost:8080/users/42/profile
	// users.Get("/messages/{id:uint64}", userMessageHandler) // http://localhost:8080/users/messages/1

	// users := app.Party("/users", myAuthMiddlewareHandler)
	// 我理解 这里用括号包起来 只是 为了好看，将同属一组的路由包在一起
	// {
	// 	users.Get("/{id:uint64}/profile", userProfileHandler)  // http://localhost:8080/users/42/profile
	// 	users.Get("/messages/{id:uint64}", userMessageHandler) // http://localhost:8080/users/messages/1
	// }

	// 或者使用 PartyFunc()方法接受子路由
	app.PartyFunc("/users", func(users iris.Party) {
		// Use appends Handler(s) to the current Party's routes and child routes. If the current Party is the root, then it registers the middleware to all child Parties' routes too.
		users.Use(myAuthMiddlewareHandler)      // Use() 在当前分组路由和子路由 增加一个处理中间件。如果当前分组是根路由，它会注册到所有子分组路由上
		users.Get("/", func(ctx iris.Context) { // https://iris-go.com/start/
			ctx.HTML("<h1>User Index</h1>")
		})
		users.Get("/{id:uint64}/profile", userProfileHandler)  // http://localhost:8080/users/42/profile
		users.Get("/messages/{id:uint64}", userMessageHandler) // http://localhost:8080/users/messages/1
	})

	app.Run(iris.Addr(":8080"))
}

func myAuthMiddlewareHandler(ctx iris.Context) {
	// 验证会员权限等
	ctx.HTML("<h1>User Auth</h1>")
	ctx.Next() // 处理下一个中间件
}

func userProfileHandler(ctx iris.Context) {
	ctx.HTML("<h1>User Profile</h1>")
}

func userMessageHandler(ctx iris.Context) {
	ctx.HTML("<h1>User Message</h1>")
}
