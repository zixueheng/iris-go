package main

import (
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"

	"github.com/kataras/iris/v12/mvc"

	"github.com/kataras/iris/v12/sessions"
)

func main() {
	app := iris.New()

	app.Use(recover.New())
	app.Logger().SetLevel("debug")

	// 创建并配置一个控制器（基于Party使用的路由，并使用 basicMVC 函数配置这个控制器）
	mvc.Configure(app.Party("/basic"), basicMVC)
	// 和下面的是 一样的效果
	// mvc.New(app.Party("/basic").Configure(basicMVC)

	app.Run(iris.Addr(":8080"))
}

// basicMVC 配置函数
func basicMVC(app *mvc.Application) {
	// You can use normal middlewares at MVC apps of course.
	// 在当前分组上 注册一个 处理中间件
	app.Router.Use(func(ctx iris.Context) {
		ctx.Application().Logger().Infof("Path: %s", ctx.Path())
		ctx.Next() // 记得要 传给下一个中间件
	})

	// Register dependencies which will be binding to the controller(s),
	// can be either a function which accepts an iris.Context and returns a single value (dynamic binding)
	// or a static struct value (service).
	app.Register( // 通过依赖注入一些变量 到 控制器中（即 结构体 basicController 中的两个字段）
		sessions.New(sessions.Config{}).Start, // 注入 session管理器的 初始化 方法体，方法体（Start()）的返回值 会绑定到 结构体basicController的 Session 字段
		&prefixedLogger{prefix: "DEV"},        // prefixedLogger 结构体已实现了 LoggerService接口的方法，所以 这里他能 绑定到 basicController 的 Logger(LoggerService类型) 字段
	)

	// GET: http://localhost:8080/basic
	// GET: http://localhost:8080/basic/custom
	// GET: http://localhost:8080/basic/custom2
	// Handel() 可以将 结构体参数上的所有方法转化成 子路由
	// 如果参数结构体 有 BeforeActivation 和 AfterActivation 将在 控制器激活前后调用
	app.Handle(new(basicController))

	// All dependencies of the parent *mvc.Application
	// are cloned to this new child,
	// thefore it has access to the same session as well.
	// GET: http://localhost:8080/basic/sub
	app.Party("/sub"). // 返回一个子控制器（路径 baisc/sub ），当前控制器的依赖变量会自动传递给子控制器
				Handle(new(basicSubController)) // 使用 basicSubController 结构体的方法做路由
}

// If controller's fields (or even its functions) expecting an interface
// but a struct value is binded then it will check
// if that struct value implements
// the interface and if true then it will add this to the
// available bindings, as expected, before the server ran of course,
// remember? Iris always uses the best possible way to reduce load
// on serving web resources.

// LoggerService 日志器
type LoggerService interface {
	Log(string)
}

type prefixedLogger struct {
	prefix string
}

func (s *prefixedLogger) Log(msg string) {
	fmt.Printf("%s: %s\n", s.prefix, msg)
}

// 当前控制器 结构体
type basicController struct {
	Logger LoggerService

	Session *sessions.Session
}

// 控制器激活前 做一些事情
func (c *basicController) BeforeActivation(b mvc.BeforeActivation) {
	// b.Dependencies().Add/Remove
	// b.Router().Use/UseGlobal/Done
	// 和已知的任何标准 API 调用

	// 这个不执行，不知道什么问题
	// b.Handle("Get", "/hello", "Custom", func(ctx iris.Context) {
	// 	// ctx.Text("Hello")
	// 	ctx.Next()
	// })

	// HandleMany 和 上面 app.Handle() 一样
	// 只是可以注册多个路由（用空格分开），第三个参数是 是要调用的方法名，第四个是前置中间所以要 Next() 执行下一个
	b.HandleMany("GET", "/custom /custom2", "Custom", func(ctx iris.Context) {
		ctx.Application().Logger().Info("Custom Log")
		ctx.Next()
	})
}

func (c *basicController) AfterActivation(a mvc.AfterActivation) {
	if a.Singleton() {
		panic("basicController should be stateless, a request-scoped, we have a 'Session' which depends on the context.")
	}
}

// 对应 http://localhost:8080/basic
func (c *basicController) Get() string {
	count := c.Session.Increment("count", 1)

	body := fmt.Sprintf("Hello from basicController\nTotal visits from you: %d", count)
	c.Logger.Log(body)
	return body
}

func (c *basicController) Custom() string {
	return "custom"
}

type basicSubController struct {
	Session *sessions.Session // 这里的会话是从父控制器复制过来
}

func (c *basicSubController) Get() string {
	count := c.Session.GetIntDefault("count", 1)
	// count := c.Session.Increment("count", 1) //  session 和 父控制器的session 是同一个
	return fmt.Sprintf("Hello from basicSubController.\nRead-only visits count: %d", count)
}
