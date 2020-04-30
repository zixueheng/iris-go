package main

import (
	"io"
	"os"
	"regexp"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"

	"iris-go/another"

	"github.com/kataras/iris/v12/middleware/logger"
)

func main() {
	another.HelloWorld() // 引用自定义包，go mod 模式下 import 要加上 module 名 iris-go

	app := iris.New() // 返回全新的 *iris.Application 实例
	app.Logger().SetLevel("debug")

	file := newLogFile()
	defer file.Close()
	// app.Logger().SetOutput(file) // 这里可以设置日志输出的地方 文件 或控制台 或 其他的地方
	app.Logger().SetOutput(io.MultiWriter(file, os.Stdout)) // 也可以同时输出到多个地方

	app.Use(recover.New()) // Recover 会从paincs中恢复并返回 500 错误码
	// 使用日志中间件
	app.Use(logger.New(logger.Config{
		// Status displays status code
		Status: true,
		// IP displays request's remote address
		IP: true,
		// Method displays the http method
		Method: true,
		// Path displays the request path
		Path: true,
		// Query appends the url query to the Path.
		Query: true,

		// if !empty then its contents derives from `ctx.Values().Get("logger_message")
		// will be added to the logs.
		MessageContextKeys: []string{"logger_message"},

		// if !empty then its contents derives from `ctx.GetHeader("User-Agent")
		MessageHeaderKeys: []string{"User-Agent"},
	}))

	// 也可以使用 Default() 返回的是一个 已包含 Logger 和 Recovery 等中间件的 *iris.Application 实例
	// 并已注册 ./views 模板 ./locales 语言包
	// app := iris.Default()

	// 设置静态文件目录 如图片 JS CSS 等
	// app.HandleDir("/static", "./assets", iris.DirOptions{ShowList: true, Gzip: true})

	// Resource: http://localhost:8080
	app.Handle("GET", "/", func(ctx iris.Context) {
		ctx.Application().Logger().Infof("请求路径: %s", ctx.Path()) // 示例：手动输出日志
		ctx.HTML("<h1>Welcome111</h1>")
	})

	// same as app.Handle("GET", "/ping", [...])
	// Resource: http://localhost:8080/ping
	app.Get("/ping", func(ctx iris.Context) {
		ctx.WriteString("pong")
	})

	// Resource: http://localhost:8080/hello
	app.Get("/hello", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"message": "Hello Iris!"})
	})

	// 多级路径匹配
	// http://localhost:8080/assets/path1
	// http://localhost:8080/assets/path1/path2
	app.Get("/assets/{asset:path}", func(ctx iris.Context) {
		ctx.Writef("Hello from method: %s and path: %s\n", ctx.Method(), ctx.Path())
	})

	// 触发一个 panic
	app.Get("/panichere", func(ctx iris.Context) {
		panic("a panic here")
	})

	// 多参数路由
	// 如果不写参数类型 则默认 string，如下面的 {name:string} 和 {name} 是一样的效果
	// http://localhost:8080/users/100/jack
	app.Get("/users/{id:uint64}/{name:string}", func(ctx iris.Context) {
		id := ctx.Params().GetUint64Default("id", 0)
		name := ctx.Params().GetStringDefault("name", "")
		ctx.Text("Get Id: %d, Name: %s", id, name)
	})

	// 注册支持所有 HTTP 方法(get post put delelte ...)的路由
	app.Any("/any", func(ctx iris.Context) {
		ctx.Writef("Hello from method: %s and path: %s\n", ctx.Method(), ctx.Path())
	})

	// 自定义路由方式：
	// The RegisterFunc can accept any function that returns a func(paramValue string) bool. Or just a func(string) bool.
	// If the validation fails then it will fire 404 or whatever status code the else keyword has.
	// RegisterFun 函数可以接受一个返回 func(paramValue string) bool 的函数 或者直接一个 a func(string) bool 函数。如果验证失败将返回404或者 在 else 其他的状态码

	// 1、正则表达式匹配路由
	phoneRegex, _ := regexp.Compile("^1[0-9]{10}$")
	// Register your custom argument-less macro function to the :string param type.
	// MatchString is a type of func(string) bool, so we use it as it is.
	app.Macros().Get("string").RegisterFunc("phone", phoneRegex.MatchString) // 注册一个无参数的宏函数 匹配 :string 参数类型
	// 使用如下
	// http://localhost:8080/cellphone/15215657185
	app.Get("/cellphone/{phone:string phone()}", func(ctx iris.Context) {
		ctx.Writef("Phone: %s", ctx.Params().Get("phone")) // 结果写入 response 中
	})
	// 正则还有一个简单的方式（不过测试没效果，还不知道怎么回事）
	app.Get("/contact{phone:string regexp(^1[0-9]{10}$)}", func(ctx iris.Context) {
		ctx.Writef("Phone: %s", ctx.Params().Get("phone")) // 结果写入 response 中
	})

	// 2、注册接受两个 int 参数的 宏函数
	app.Macros().Get("string").RegisterFunc("range", func(minLength, maxLength int) func(string) bool {
		return func(paramValue string) bool {
			return len(paramValue) >= minLength && len(paramValue) <= maxLength
		}
	})
	// 使用如下:（参数3和10 传给给 minLength和maxLength）else 400 意思是否则返回 400 Bad Request 而不返回默认的 404 Not Found
	app.Get("/limitchar/{name:string range(3,10) else 400}", func(ctx iris.Context) {
		name := ctx.Params().Get("name")
		ctx.Writef(`Hello %s | the name should be between 3 and 10 characters length
		otherwise this handler will not be executed`, name)
	})

	// 3、注册接受字符串切片的 宏函数
	app.Macros().Get("string").RegisterFunc("has", func(validNames []string) func(string) bool {
		return func(paramValue string) bool {
			for _, validName := range validNames {
				if validName == paramValue {
					return true
				}
			}
			return false
		}
	})
	// 使用如下:传入的参数 必须在 限定的几项中
	app.Get("/static_validation/{name:string has([kataras,gerasimos,maropoulos])}", func(ctx iris.Context) {
		name := ctx.Params().Get("name")
		ctx.Writef(`Hello %s | the name should be "kataras" or "gerasimos" or "maropoulos"
		otherwise this handler will not be executed`, name)
	})

	// 配置主机
	app.ConfigureHost(func(h *iris.Supervisor) { // 监控服务器
		h.RegisterOnShutdown(func() { // 注册一个服务器关闭的回调 打印
			println("server terminated")
		})
	})

	// Iris 的默认行为是接受和注册像 /api/user 这类路径的路由，而且路径末端不带斜杠。
	// 如果客户端尝试访问 $your_host/api/user/ ，则 Iris 路由器将自动将其重定向到 $your_host/api/user
	// 如果想为 /api/user 和 /api/user/ 路径，保持相同的处理程序和路由而无需重定向 (常见情况) ，那么仅仅使用 iris.WithoutPathCorrectionRedirection 选项即
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed), iris.WithoutPathCorrectionRedirection)
}

// 生成一个日志的 txt 文件名
func todayFilename() string {
	today := time.Now().Format("2006-01-02")
	return today + ".txt"
}

// 创建文件并，返回文件句柄
func newLogFile() *os.File {
	filename := todayFilename()
	// Open the file, this will append to the today's file if server restarted.
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	return f
}
