package main

import (
	"github.com/kataras/iris/v12"

	"github.com/kataras/iris/v12/sessions"
)

// 有时需要在同一用户的请求之间临时存储数据，例如在提交表单后出现错误或成功消息。Iris Sessions 程序包支持闪存消息。

// The Session contains the following methods to store, retrieve and remove flash messages.

// SetFlash(key string, value interface{})
// HasFlash() bool
// GetFlashes() map[string]interface{}

// PeekFlash(key string) interface{}
// GetFlash(key string) interface{}
// GetFlashString(key string) string
// GetFlashStringDefault(key string, defaultValue string) string

// DeleteFlash(key string)
// ClearFlashes()

// The method names are self-explained. For example, if you need to get a message and remove it at the next request use the GetFlash.
// When you only need to retrieve but don't remove then use the PeekFlash.

// Flash messages are not stored in a database.

func main() {
	app := iris.New()
	sess := sessions.New(sessions.Config{Cookie: "myappsessionid", AllowReclaim: true})

	app.Get("/set", func(ctx iris.Context) {
		s := sess.Start(ctx)
		s.SetFlash("name", "iris")
		ctx.Writef("Message set, is available for the next request")
	})

	app.Get("/get", func(ctx iris.Context) {
		s := sess.Start(ctx)
		name := s.GetFlashString("name") // 取出 name 并再下次请求时删除（意思这里取过了 下次就取不到了）
		if name == "" {
			ctx.Writef("Empty name!!")
			return
		}
		ctx.Writef("Hello %s", name)
	})

	app.Get("/test", func(ctx iris.Context) {
		s := sess.Start(ctx)
		name := s.GetFlashString("name")
		if name == "" {
			ctx.Writef("Empty name!!")
			return
		}

		ctx.Writef("Ok you are coming from /set ,the value of the name is %s", name)
		ctx.Writef(", and again from the same context: %s", name)
	})

	app.Listen(":8080")
}
