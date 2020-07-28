# go-serverchan

+ An unofficial serverchan client for golang

### ServerChan

+ Please see [Server酱](http://sc.ftqq.com/3.version)
+ Attention:

> 每人每天发送上限500条，相同内容5分钟内不能重复发送，不同内容10秒内不能连续发送 不同内容一分钟只能发送30条。
> 主要是防止程序出错的情况。注意，因为之前频繁的死循环程序导致费用飙升，现在每天调用接口超过1000次的用户将被系统自动拉黑。

### Usage

```go
package main

import (
    serverchan "github.com/Aoi-hosizora/go-serverchan"
    "log"
)

func main() {
    client := serverchan.NewClient()
    _, _ = client.Send("sckey", "title", "message")

    // logger
    client.SetLogger(&MyLogger{})
}

type MyLogger struct{}

func (m *MyLogger) Log(sckey string, title string, err error) {
    // xxx
}
```
