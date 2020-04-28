## 如果你想将静态文件嵌入可执行文件内部以不依赖于系统目录，则可以使用 go-bindata 将文件转换为程序内的 []byte

```shell
# 下载工具
$ go get -u github.com/go-bindata/go-bindata/...
# 编码文件夹
$ go-bindata ./assets/...
# 编译项目
$ go build
# 运行可执行文件
$ ./embedding-files-into-app
```
assets 文件夹里面的文件 运行的时候不会使用，可以删除

如果要 gzip 压缩的方式 需要另外一个 工具