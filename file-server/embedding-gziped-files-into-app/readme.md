## 编码并gzip 静态文件

```shell
# 下载工具，不过报错了 不知道什么问题，使用不了
$ go get -u github.com/kataras/bindata/cmd/bindata
$ bindata ./assets/...
$ go build
$ ./embedding-gziped-files-into-app
```
"physical" files are not used, you can delete the "assets" folder and run the example.