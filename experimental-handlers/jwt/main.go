// iris provides some basic middleware, most for your learning curve.
// You can use any net/http compatible middleware with iris.FromStd wrapper.
//
// JWT net/http video tutorial for golang newcomers: https://www.youtube.com/watch?v=dgJFeqeXVKw
//
// This middleware is the only one cloned from external source: https://github.com/auth0/go-jwt-middleware
// (because it used "context" to define the user but we don't need that so a simple iris.FromStd wouldn't work as expected.)
package main

import (
	"github.com/kataras/iris/v12"

	"github.com/iris-contrib/middleware/jwt"
)

const secretKey = "SecretKey"

func getTokenHandler(ctx iris.Context) {
	// 获取一个 Token，参数一：签名方法、参数二：要保存的数据
	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":      "1",
		"username": "jack",
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, _ := token.SignedString([]byte(secretKey))

	ctx.HTML(`Token: ` + tokenString + `<br/><br/>
    <a href="/secured?token=` + tokenString + `">/secured?token=` + tokenString + `</a>`)
}

func myAuthenticatedHandler(ctx iris.Context) {
	user := ctx.Values().Get("jwt").(*jwt.Token)

	ctx.Writef("This is an authenticated request\n")
	ctx.Writef("Claim content:\n")

	data := user.Claims.(jwt.MapClaims)
	for key, value := range data {
		ctx.Writef("%s = %s\n", key, value)
	}
}

func main() {
	app := iris.New()

	// 构建一个 jwt 实例
	j := jwt.New(jwt.Config{
		// Extract by "token" url parameter.
		// Extractor: jwt.FromParameter("token"), // 从请求参数中提起 token

		// 一般情况 是提取 header 中的 Authorization
		// Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIifQ.4aQan-XvJBjDUDbCVnuh2P_xy54b2aRKKsgKcHUa8uw
		Extractor: jwt.FromAuthHeader,

		// SecretKey
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	app.Get("/", getTokenHandler)

	// j.Serve 验证令牌
	app.Get("/secured", j.Serve, myAuthenticatedHandler)
	app.Listen(":8080")
}
