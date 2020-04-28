package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/hero"
)

// 依赖注入 示例

func main() {

	app := iris.New()

	// 1. Path Parameters - Built-in Dependencies 内建依赖
	// 1、根据路径参数进行注入，{to:string}，与 func hello(to string)，进行注入
	helloHandler := hero.Handler(hello)
	app.Get("/{to:string}", helloHandler) // http://localhost:8080/jack

	// 2. Services - Static Dependencies 静态依赖
	// 2、注册一个结构体变量 &myTestService 用于依赖， 变量将注入到 helloService()的第二个参数 service（类型Service）
	hero.Register(&myTestService{ // 注意：&myTestService 已实现了 Service 接口，所以才能传给 service
		prefix: "Service: Hello",
	})

	helloServiceHandler := hero.Handler(helloService)
	app.Get("/service/{to:string}", helloServiceHandler) // to 传递给 helloService()第一个参数，helloService()的第二个参数是注入的 &myTestService

	// 3. Per-Request - Dynamic Dependencies 动态依赖
	hero.Register(func(ctx iris.Context) (form LoginForm) {
		// it binds the "form" with a
		// x-www-form-urlencoded form data and returns it.
		ctx.ReadForm(&form) // 要求 参数 结构体 字段的 tag 为 form
		// ReadForm binds the formObject with the form data it supports any kind of type, including custom structs.
		// It will return nothing if request data are empty. The struct field tag is "form".
		return
	})

	loginHandler := hero.Handler(login)
	app.Post("/login", loginHandler) // // 可以用 Postman Post 模拟

	// http://localhost:8080/your_name
	// http://localhost:8080/service/your_name
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}

func hello(to string) string {
	return "Hello " + to
}

// Service 接口
type Service interface {
	SayHello(to string) string
}

type myTestService struct {
	prefix string
}

// 结构体实现 Service 接口
func (s *myTestService) SayHello(to string) string {
	return s.prefix + " " + to
}

func helloService(to string, service Service) string {
	return service.SayHello(to)
}

// LoginForm 结构体
type LoginForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func login(form LoginForm) string {
	return "Hello " + form.Username
}
