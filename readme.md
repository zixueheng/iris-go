# Iris Golang 学习

自动重启工具（监控源码改变自动重启进程）https://github.com/kataras/rizla

`go get -u github.com/kataras/rizla`

```shell
$ rizla main.go #single project monitoring
$ rizla C:/myprojects/project1/main.go C:/myprojects/project2/main.go #multi projects monitoring
$ rizla -walk main.go #prepend '-walk' only when the default file changes scanning method doesn't works for you.
$ rizla -delay=5s main.go # if delay > 0 then it delays the reload, also note that it accepts the first change but the rest of changes every "delay".
```

## 路由定义
https://github.com/kataras/iris/wiki/Routing-path-parameter-types

## Golang交叉编译平台的二进制文件
```shell
# mac上编译linux和windows二进制
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build 
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build 
 
# linux上编译mac和windows二进制
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build 
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build
 
# windows上编译mac和linux二进制
SET CGO_ENABLED=0 SET GOOS=darwin SET GOARCH=amd64 go build main.go
SET CGO_ENABLED=0 SET GOOS=linux SET GOARCH=amd64 go build main.go
```

### linux 上运行二进制文件
```shell
# 修改权限命令
chmod 777 程序名称

# 后台运行的命令
nohup ./程序名 & 

# 不输出错误信息
nohup ./程序名 >/dev/null 2>&1 &

# 如果要关闭程序，可以使用命令”ps” 查看后台程序的pid，然后使用“kill 程序pid”命令，关闭程序比如程序名为test，可以用如下命令查询
ps aux|grep test
```

#### 命令 nohup 和 & 
&：指在后台运行
nohup：不挂断的运行，注意并没有后台运行的功能，，就是指，用nohup运行命令可以使命令永久的执行下去，和用户终端没有关系，例如我们断开SSH连接都不会影响他的运行，注意了nohup没有后台运行的意思；&才是后台运行

&：是指在后台运行，但当用户退出(挂起)的时候，命令自动也跟着退出
那么，我们可以巧妙的吧他们结合起来用就是
nohup COMMAND &
这样就能使命令永久的在后台执行
例如：
1. sh test.sh &  
将sh test.sh任务放到后台 ，即使关闭xshell退出当前session依然继续运行，但标准输出和标准错误信息会丢失（缺少的日志的输出）

将sh test.sh任务放到后台 ，关闭xshell，对应的任务也跟着停止。
2. nohup sh test.sh  
将sh test.sh任务放到后台，关闭标准输入，终端不再能够接收任何输入（标准输入），重定向标准输出和标准错误到当前目录下的nohup.out文件，即使关闭xshell退出当前session依然继续运行。
3. nohup sh test.sh  & 
将sh test.sh任务放到后台，但是依然可以使用标准输入，终端能够接收任何输入，重定向标准输出和标准错误到当前目录下的nohup.out文件，即使关闭xshell退出当前session依然继续运行