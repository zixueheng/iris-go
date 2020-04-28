package main

import (
	"github.com/kataras/iris/v12"

	"github.com/kataras/iris/v12/sessions"
)

var (
	cookieNameForSessionID = "mycookiesessionnameid"
	// 创建一个 session 管理器
	sess = sessions.New(sessions.Config{Cookie: cookieNameForSessionID, AllowReclaim: true})
)

func secret(ctx iris.Context) {
	// Check if user is authenticated
	if auth, _ := sess.Start(ctx).GetBoolean("authenticated"); !auth {
		ctx.StatusCode(iris.StatusForbidden)
		return
	}

	// Print secret message
	ctx.WriteString("The cake is a lie!")
}

func login(ctx iris.Context) {
	session := sess.Start(ctx) // 从 ctx 中获取 会话，如果没有就创建一个

	// Authentication goes here
	// ...

	// Set user as authenticated
	session.Set("authenticated", true)

}

func logout(ctx iris.Context) {
	session := sess.Start(ctx)

	// Revoke users authentication
	session.Set("authenticated", false)
	// 或者 直接销毁
	// session.Destroy() // 销毁服务端 会话，但是客户端的 cookie 还是存在的

	sess.Destroy(ctx) // 销毁服务端会话，也删除客户端的 cookie
}

func main() {
	app := iris.New()

	app.Get("/secret", secret)
	app.Get("/login", login)
	app.Get("/logout", logout)

	app.Listen(":8080")
}
