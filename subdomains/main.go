package main

import "github.com/kataras/iris/v12"

// 子域名

// windows 操作系统中打开 C:\Windows\System32\Drivers\etc\hosts 文件并且加入一下内容：
// 127.0.0.1 mydomain.com
// 127.0.0.1 admin.mydomain.com

func main() {
	app := iris.New()

	// subdomains works with all available routers, like other features too.

	// no order, you can register subdomains at the end also.
	admin := app.Party("admin.") // 创建一个子域名 分组，然后在这个分组下的子路由都是这个子域名下的
	// 也可以使用下面的代码
	// admin := app.Subdomain("admin")
	{
		// admin.mydomain.com
		admin.Get("/", func(c iris.Context) {
			c.Writef("INDEX FROM admin.mydomain.com")
			// 更多方法演示
			method := c.Method()
			subdomain := c.Subdomain()
			path := c.Path()
			www := c.IsWWW()

			c.Writef("\nIs www: %v \n\n", www)
			c.Writef("Method: %s\nSubdomain: %s\nPath: %s", method, subdomain, path)
		})
		// admin.mydomain.com/hey
		admin.Get("/hey", func(c iris.Context) {
			c.Writef("HEY FROM admin.mydomain.com/hey")
		})
		// admin.mydomain.com/hey2
		admin.Get("/hey2", func(c iris.Context) {
			c.Writef("HEY SECOND FROM admin.mydomain.com/hey")
		})
	}

	// 下面代码是演示 多个 子域名 在同一 app 中
	dashboard := app.Party("dashboard.")
	{
		dashboard.Get("/", func(ctx iris.Context) {
			ctx.Writef("HEY FROM dashboard")
		})
	}
	system := app.Party("system.")
	{
		system.Get("/", func(ctx iris.Context) {
			ctx.Writef("HEY FROM system")
		})
	}

	// mydomain.com/
	app.Get("/", func(c iris.Context) {
		c.Writef("INDEX FROM no-subdomain hey")
	})

	// mydomain.com/hey
	app.Get("/hey", func(c iris.Context) {
		c.Writef("HEY FROM no-subdomain hey")
	})

	// http://admin.mydomain.com
	// http://admin.mydomain.com/hey
	// http://admin.mydomain.com/hey2
	// http://mydomain.com
	// http://mydomain.com/hey
	app.Run(iris.Addr("mydomain.com:80")) // for beginners: look ../hosts file
}
