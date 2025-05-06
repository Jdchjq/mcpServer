## MCP SERVER

这是基于 mcp-go 实现的一些 mcp server

## server

### 细节

目前实现的功能有

- weather 天气相关的功能，用于查询天气数据。底层使用的是和风天气 API

### 一些补充

目前 mcp-go 没有实现对返回结果进行注释，本项目通过重写自定义类型的 UnmarshalJSON 方法，将注释的 tag 赋值到每个 description 字段中，方便 AI 在获取结果后能充分、正确的解读（一定程度增加了 token 量）

## client

### 使用方式

执行输入参数：  
1、configDir  
表示配置文件所在目录，路径下的配置文件需要 config/config.yaml 、config/private_key.pem  
因为配置比较多，没有改成环境变量的输入方式  
2、transport  
mcp server 的运行方式，sse 表示采用 http 推流方式运行，stdio 表示采用进程间通信方式运行

客户端配置  
1、通过 stdio（进程间通信）方式调用 server  
先使用 go build -0 ./bin/weather ./cmd/weather 将程序编译成可执行文件  
最后在 command 填入可执行文件的路径，env 配置 go 的相关环境变量

```json
"my-weather-server": {
  "command": "/Users/jundongchen/Documents/go/src/smart-customer/bin/weather",
  "args": [
    "--configDir=/Users/jundongchen/Documents/go/src/smart-customer",
    "--transport=stdio"
  ],
  "env": {
    "GOPATH": "/Users/jundongchen/Documents/go",
    "GOMODCACHE": "/Users/jundongchen/Documents/go/pkg/mod"
  }
}

```

2、通过 sse 方式调用 server  
此方式适用于需要远程调用 mcp server, 在使用前，需要使用`go run ./cmd/weather/main.go --configDir=/Users/jundongchen/Documents/go/src/smart-customer -t=sse` 运行 mcp server  
客户端 mcp json 配置如下：

```json
"my-weather-server-sse": {
  "url": "http://localhost:8080/sse",
  "args": [
    "--configDir=/Users/jundongchen/Documents/go/src/smart-customer",
    "--transport=stdio"
  ],
  "env": {
    "GOPATH": "/Users/jundongchen/Documents/go",
    "GOMODCACHE": "/Users/jundongchen/Documents/go/pkg/mod"
  },
  "transportType": "sse",
},

```
