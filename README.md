# fingerprintx 描述

当前为`https://github.com/expzhizhuo/fingerprintx`仓库的克隆，由于官方支持的tcp协议过少，特此克隆仓库进行自定义的协议添加补充。

## 编写插件

插件位置位于`pkg/plugins/services/`下，新建对应协议的目录参考其他的协议编写即可。

写完插件需要在`pkg/scan/plugin_list.go`文件中引入。在`pkg/plugins/types.go`中声明类型，否在插件会报错。
例如：
```go

type ServiceMongoDB struct {
	PacketType   string `json:"packetType"` // 服务器返回的数据包类型（例如：握手或错误）
	PacketData   string `json:"PacketData"` // 响应数据
	ErrorMessage string `json:"errorMsg"`   // 如果服务器返回错误数据包，则为错误消息
	ErrorCode    int    `json:"errorCode"`  // 如果服务器返回错误数据包，则为错误代码
}
```
插件的写法也需要注意，这里如果是需要协议通信，我们可以使用WireShare去抓hex string，然后可以使用下面函数转换成byte类型
```go
func createBuildInfoCommand() []byte {
	hexStr := "hex内容"
	cmd := make([]byte, len(hexStr)/2)
	for i := 0; i < len(hexStr); i += 2 {
		decodedByte, _ := hex.DecodeString(hexStr[i : i+2])
		cmd[i/2] = decodedByte[0]
	}
	return cmd
}
```
其中函数`Priority()`是配置插件的优先级，例如2表示这个插件的优先级为2，如果有的协议可以通过http/https访问的话，个人建议将优先级设置的靠前一些，可以有效的避免识别成http
```go
func (p *MongoDBPlugin) Priority() int {
	return 2
}
```

其中函数`PortPriority()`是配置其默认端口
```go
func (p *MongoDBPlugin) PortPriority(port uint16) bool {
	return port == 27017
}
```

## 最后

如果插件的写的有错误请提交issues。如果你有自己编写好的插件可以通过pr和本仓库合并。

## 明谢 

- https://github.com/expzhizhuo/fingerprintx