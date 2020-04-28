package main

import (
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New() // defaults to these

	tmpl := iris.HTML("./templates", ".html")
	tmpl.Reload(true) // 对每一个请求 重新加载 模板 (开发模式使用)
	// default template funcs are:
	//
	// - {{ urlpath "mynamedroute" "pathParameter_ifneeded" }}
	// - {{ render "header.html" }}
	// - {{ render_r "header.html" }} // partial relative path to current page
	// - {{ yield }}
	// - {{ current }}
	tmpl.AddFunc("greet", func(s string) string {
		return "Greetings " + s + "!"
	})
	app.RegisterView(tmpl)

	app.Get("/", hi)

	about := app.Get("/about", func(ctx iris.Context) {
		ctx.HTML("<h1>About Page</h1>")
	})
	about.Name = "About"

	// http://localhost:8080
	app.Listen(":8080", iris.WithCharset("utf-8")) // defaults to that but you can change it.
}

func hi(ctx iris.Context) {
	ctx.ViewData("Title", "Hi Page")
	ctx.ViewData("Name", "iris") // {{.Name}} will render: iris
	// ctx.ViewData("", myCcustomStruct{})
	ctx.ViewData("Url", "About") // 关于 URL 的 使用参考 template_route
	ctx.View("hi.html")
}

/*
Note:
In case you're wondering, the code behind the view engines derives from the "github.com/kataras/iris/v12/view" package,
access to the engines' variables can be granded by "github.com/kataras/iris/v12" package too.
    iris.HTML(...) is a shortcut of view.HTML(...)
    iris.Django(...)     >> >>      view.Django(...)
    iris.Pug(...)        >> >>      view.Pug(...)
    iris.Handlebars(...) >> >>      view.Handlebars(...)
    iris.Amber(...)      >> >>      view.Amber(...)
*/
