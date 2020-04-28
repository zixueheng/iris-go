package main

import (
	"github.com/kataras/iris/v12"
)

//  查询参数 和 表单示例

func main() {
	app := iris.Default()

	// 1 查询参数
	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe.
	app.Get("/welcome", func(ctx iris.Context) {
		firstname := ctx.URLParamDefault("firstname", "Guest")

		// shortcut for ctx.Request().URL.Query().Get("lastname").
		lastname := ctx.URLParam("lastname")

		ctx.Writef("Hello %s %s", firstname, lastname) // Hello Jane Doe
	})

	// 2 表单
	// Multipart/Urlencoded Form
	app.Post("/form_post", func(ctx iris.Context) {
		message := ctx.FormValue("message")               // 取表单参数 message
		nick := ctx.FormValueDefault("nick", "anonymous") // 取表单参数 nick ，没有默认 anonymous

		// 返回 json 数据
		ctx.JSON(iris.Map{
			"status":  "posted",
			"message": message,
			"nick":    nick,
		})
	})

	// 3 带查询参数的表单
	// POST /post?id=1234&page=1 HTTP/1.1
	// Content-Type: application/x-www-form-urlencoded
	app.Post("/post", func(ctx iris.Context) {
		id := ctx.URLParam("id")
		page := ctx.URLParamDefault("page", "0")
		name := ctx.FormValue("name")
		message := ctx.FormValue("message")
		// or `ctx.PostValue` for POST, PUT & PATCH-only HTTP Methods.

		// 控制台 打印日志
		app.Logger().Infof("id: %s; page: %s; name: %s; message: %s", id, page, name, message) // id: 1234; page: 1; name: manu; message: this_is_great
	})

	app.Run(iris.Addr(":8080"))
}
