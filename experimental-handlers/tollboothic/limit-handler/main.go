package main

import (
	"github.com/kataras/iris/v12"

	"github.com/didip/tollbooth"
	"github.com/iris-contrib/middleware/tollboothic"
)

// $ go get github.com/didip/tollbooth
// $ go run main.go

// 请求频率限制 中间件
// 超出频率 返回错误码 429 Too Many Requests

func main() {
	app := iris.New()

	limiter := tollbooth.NewLimiter(1, nil) //  这里设置的1秒内 最多 1个请求
	//
	// or create a limiter with expirable token buckets
	// This setting means:
	// create a 1 request/second limiter and
	// every token bucket in it will expire 1 hour after it was initially set.
	// limiter := tollbooth.NewLimiter(1, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})

	app.Get("/", tollboothic.LimitHandler(limiter), func(ctx iris.Context) {
		ctx.HTML("<b>Hello, world!</b>")
	})

	app.Listen(":8080")
}

// Read more at: https://github.com/didip/tollbooth
