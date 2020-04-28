## Iris 为 httpexpect 提供了令人难以置信的支持，它是一个 Web 应用程序测试框架。  iris/httptest 子包为 Iris + httpexpect 体统了辅助方法。

httpexpect地址：https://github.com/gavv/httpexpect

## 如果你喜欢 Go 的标准 net/http/httptest 包，你仍然可以使用它。

### 第一个例子 使用 iris/httptest 来测试基本身份验证

main.go 和 main_test.go
`$ go test -v`

### 第二个例子 Cookies
`$ go test -v -run=TestCookiesBasic$`