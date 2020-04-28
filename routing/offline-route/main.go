package main

import (
	"github.com/kataras/iris/v12"
)

// 在 Iris 中有一种特别的方法。 它被称为 None 并且可以用它对外部隐藏一个路由， 但是你却可以通过 Context.Exec 方法从其他路由处理程序中调用它。
// 每个 API 处理方法都会返回路由的 值。路由的 IsOnline 方法报告该路由的当前状态。
// 你可以通过路由 Route.Method 字段的值，将路由的状态从离线改变在线。
// 当然，在服务时路由器的每次更改都需要一个 app.RefreshRouter() 调用，该调用可以安全使用

func main() {
	app := iris.New()

	// 不可见的路由，即通过外部不可访问
	none := app.None("/invisible/{username}", func(ctx iris.Context) {
		ctx.Writef("Hello %s with method: %s", ctx.Params().Get("username"), ctx.Method())

		if from := ctx.Values().GetString("from"); from != "" {
			ctx.Writef("\nI see that you're coming from %s", from)
		}
	})

	// 更改 none 路由 的状态
	app.Get("/change", func(ctx iris.Context) {

		if none.IsOnline() { // 如果 none 在线
			none.Method = iris.MethodNone // 将 none 的 Method 改为 NONE
		} else {
			none.Method = iris.MethodGet // 如果不在线 则 Method 改成 Get (none := app.Get("/invisible/{username}")
		}

		// refresh re-builds the router at serve-time in order to
		// be notified for its new routes.
		app.RefreshRouter() // 路由状态改变后 要调用这个 刷新路由
	})

	// 这里定义一个其他路由 可以在内部执行 不可见路由none
	app.Get("/execute/{username}", func(ctx iris.Context) {
		username := ctx.Params().Get("username")
		if !none.IsOnline() { // none不在线
			ctx.Values().Set("from", "/execute with offline access") // 在当前 context 上设置一个 from 变量
			ctx.Exec("NONE", "/invisible/"+username)                 // 执行其他路由，这里执行路径 /invisible/... 方式 NONE，会得到其他路由的响应
			return
		}

		// same as navigating to "http://localhost:8080/invisible/iris"
		// when /change has being invoked and route state changed
		// from "offline" to "online"
		ctx.Values().Set("from", "/execute online now")
		// values and session can be
		// shared when calling Exec from a "foreign" context.
		// 	ctx.Exec("NONE", "/invisible/iris")
		// or after "/change":
		ctx.Exec("GET", "/invisible/"+username)
	})

	app.Listen(":8080")
}

// 如何运行
//     go run main.go
//     打开一个位于 http://localhost:8080/invisible/iris 的浏览器， 然后你会获得一个 404 not found 错误，
//     但是 http://localhost:8080/execute 将能够执行该路由。
//     现在，如果导航到 http://localhost:8080/change 并刷新 /invisible/iris 选项卡 ，你将会看到它。
