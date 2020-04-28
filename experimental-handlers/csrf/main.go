// This middleware provides Cross-Site Request Forgery protection.
//
// It securely generates a masked (unique-per-request) token that
// can be embedded in the HTTP response (e.g. form field or HTTP header).
// The original (unmasked) token is stored in the session, which is inaccessible
// by an attacker (provided you are using HTTPS). Subsequent requests are
// expected to include this token, which is compared against the session token.
// Requests that do not provide a matching token are served with a HTTP 403
// 'Forbidden' error response.
package main

// $ go get -u github.com/iris-contrib/middleware/...

// CSRF 跨站请求伪造（英语：Cross-site request forgery），也被称为 one-click attack 或者 session riding，通常缩写为 CSRF 或者 XSRF，
// 是一种挟制用户在当前已登录的Web应用程序上执行非本意的操作的攻击方法。
// 跟跨网站脚本（XSS）相比，XSS 利用的是用户对指定网站的信任，CSRF 利用的是网站对用户网页浏览器的信任。

// 这里主要通过 生成一个 token  嵌入到 HTTP响应中 返回给客户端，原始的 token 保存到服务器的 session中，客服端 后续的请求需要携带 token 用来和服务器session中的 token 进行验证
// 失败则返回 403

import (
	"github.com/kataras/iris/v12"

	"github.com/iris-contrib/middleware/csrf"
)

func main() {
	app := iris.New()
	app.RegisterView(iris.HTML("./views", ".html"))
	// Note that the authentication key provided should be 32 bytes
	// long and persist across application restarts.
	protect := csrf.Protect([]byte("9AB0F421E53A477C084477AEA06096F5"), // 32位的 key
		csrf.Secure(false)) // Defaults to true, but pass `false` while no https (devmode).

	users := app.Party("/user", protect)
	{
		users.Get("/signup", getSignupForm)
		// // POST requests without a valid token will return a HTTP 403 Forbidden.
		users.Post("/signup", postSignupForm)
	}

	// GET: http://localhost:8080/user/signup
	// POST: http://localhost:8080/user/signup
	app.Listen(":8080")
}

func getSignupForm(ctx iris.Context) {
	// views/user/signup.html just needs a {{ .csrfField }} template tag for
	// csrf.TemplateField to inject the CSRF token into. Easy!
	ctx.ViewData(csrf.TemplateTag, csrf.TemplateField(ctx))
	ctx.View("user/signup.html")

	// We could also retrieve the token directly from csrf.Token(ctx) and
	// set it in the request header - ctx.GetHeader("X-CSRF-Token", token)
	// This is useful if you're sending JSON to clients or a front-end JavaScript
	// framework.
}

func postSignupForm(ctx iris.Context) {
	ctx.Writef("You're welcome mate!")
}
