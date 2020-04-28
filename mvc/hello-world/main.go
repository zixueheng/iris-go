package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"

	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

// 控制器 方法 对应 路由 的解析案列

// This example is equivalent to the
// https://github.com/kataras/iris/blob/master/_examples/hello-world/main.go
//
// It seems that additional code you
// have to write doesn't worth it
// but remember that, this example
// does not make use of iris mvc features like
// the Model, Persistence or the View engine neither the Session,
// it's very simple for learning purposes,
// probably you'll never use such
// as simple controller anywhere in your app.
//
// The cost we have on this example for using MVC
// on the "/hello" path which serves JSON
// is ~2MB per 20MB throughput on my personal laptop,
// it's tolerated for the majority of the applications
// but you can choose
// what suits you best with Iris, low-level handlers: performance
// or high-level controllers: easier to maintain and smaller codebase on large applications.

// Of course you can put all these to main func, it's just a separate function
// for the main_test.go.
func newApp() *iris.Application {
	app := iris.New()
	// Optionally, add two builtin handlers
	// that can recover from any http-relative panics
	// and log the requests to the terminal.
	app.Use(recover.New())
	app.Use(logger.New())

	tmpl := iris.HTML("./views", ".html")
	tmpl.Reload(true)      // 对每一个请求 重新加载 模板 (开发模式使用)
	app.RegisterView(tmpl) // 注册模板引擎

	// Serve a controller based on the root Router, "/".
	mvc.New(app).Handle(new(ExampleController)) // New(app) 和 New(app.Party("/")) 是相等的，也就是 注册 根 / 路由
	return app
}

func main() {
	app := newApp()

	// http://localhost:8080
	// http://localhost:8080/ping
	// http://localhost:8080/hello
	// http://localhost:8080/custom_path
	app.Listen(":8080")
}

// ExampleController serves the "/", "/ping" and "/hello".
type ExampleController struct{}

// Get serves
// Method:   GET
// Resource: http://localhost:8080
func (c *ExampleController) Get() mvc.Result {
	return mvc.Response{
		ContentType: "text/html",
		Text:        "<h1>Welcome</h1>",
	}
}

// GetPing serves
// Method:   GET
// Resource: http://localhost:8080/ping
func (c *ExampleController) GetPing() string {
	return "pong"
}

// GetHello serves
// Method:   GET
// Resource: http://localhost:8080/hello
func (c *ExampleController) GetHello() interface{} {
	return map[string]string{"message": "Hello Iris!"} // 会根据返回值类型 自动转化 这里 map 会返回给客户端 json
}

// GetHelloWorld 多大写字母
// http://localhost:8080/hello/world
func (c *ExampleController) GetHelloWorld() interface{} {
	return map[string]string{"message": "Hello World!"} // 会根据返回值类型 自动转化 这里 map 会返回给客户端 json
}

// BeforeActivation called once, before the controller adapted to the main application
// and of course before the server ran.
// After version 9 you can also add custom routes for a specific controller's methods.
// Here you can register custom method's handlers
// use the standard router with `ca.Router` to do something that you can do without mvc as well,
// and add dependencies that will be binded to a controller's fields or method function's input arguments.
// BeforeActivation 在控制器填充到 主应用之前调用（当然也在服务启动前），这里可以做一些 添加删除依赖的操作 或增加一些路由
func (c *ExampleController) BeforeActivation(b mvc.BeforeActivation) {
	anyMiddlewareHere := func(ctx iris.Context) {
		ctx.Application().Logger().Warnf("Inside /custom_path")
		ctx.Next()
	}
	b.Handle("GET", "/custom_path", "CustomHandlerWithoutFollowingTheNamingGuide", anyMiddlewareHere)

	// or even add a global middleware based on this controller's router,
	// which in this example is the root "/":
	// b.Router().Use(myMiddleware) //或者在控制器路由上 增加一些 全局中间件
}

// CustomHandlerWithoutFollowingTheNamingGuide serves
// Method:   GET
// Resource: http://localhost:8080/custom_path
func (c *ExampleController) CustomHandlerWithoutFollowingTheNamingGuide() string {
	return "hello from the custom handler without following the naming guide"
}

// GetUserBy serves
// Method:   GET
// Resource: http://localhost:8080/user/{username:string}
// By is a reserved "keyword" to tell the framework that you're going to
// bind path parameters in the function's input arguments, and it also
// helps to have "Get" and "GetBy" in the same controller.
//
// 使用 GetBy... PostBy... 参数则当做路由的参数
// 返回值 mvc.Result 是 hero.Result 类型的一个别名， 同时它也是一个接口，他的 Dispatch() 发送响应到 context 的 response writer 中
func (c *ExampleController) GetUserBy(username string) mvc.Result {
	return mvc.View{
		Name: "username.html",
		Data: username,
	}
}

// GetByWildcard 通配符路由 匹配 不是上面的 所有路由
func (c *ExampleController) GetByWildcard(path string) string {
	return path
}

/* Can use more than one, the factory will make sure
that the correct http methods are being registered for each route
for this controller, uncomment these if you want:
func (c *ExampleController) Post() {}
func (c *ExampleController) Put() {}
func (c *ExampleController) Delete() {}
func (c *ExampleController) Connect() {}
func (c *ExampleController) Head() {}
func (c *ExampleController) Patch() {}
func (c *ExampleController) Options() {}
func (c *ExampleController) Trace() {}
*/

/*
func (c *ExampleController) All() {}
//        OR
func (c *ExampleController) Any() {}
func (c *ExampleController) BeforeActivation(b mvc.BeforeActivation) {
	// 1 -> the HTTP Method
	// 2 -> the route's path
	// 3 -> this controller's method name that should be handler for that route.
	b.Handle("GET", "/mypath/{param}", "DoIt", optionalMiddlewareHere...)
}
// After activation, all dependencies are set-ed - so read only access on them
// but still possible to add custom controller or simple standard handlers.
func (c *ExampleController) AfterActivation(a mvc.AfterActivation) {}
*/
