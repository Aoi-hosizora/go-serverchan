# go-serverchan

+ An unofficial serverchan client for golang

### ServerChan

+ Please see [Server酱](http://sc.ftqq.com/3.version)

> Attention:
> 
> 每人每天发送上限500条，相同内容5分钟内不能重复发送，不同内容10秒内不能连续发送 不同内容一分钟只能发送30条。
> 主要是防止程序出错的情况。注意，因为之前频繁的死循环程序导致费用飙升，现在每天调用接口超过1000次的用户将被系统自动拉黑。

### Tips

1. Errmsg `不要重复发送同样的内容` is replaced to `duplicate message` and still with errno 1024
2. If sckey is not existed, it will return errmsg `bad pushtoken` with errno 1024 (`CheckSckey` will check this errmsg rather than errno)
3. `client.Send` will return three values, and the last value error will have 2 type: general error, server error (can be checked by `IsResponseError`)

### Usage

```go
package main

import (
    serverchan "github.com/Aoi-hosizora/go-serverchan"
)

func main() {
    client := serverchan.NewClient()
    _, _, _ = client.Send("sckey", "title", "message")
    _, _ = client.CheckSckey("sckey", "test")

    // logger
    client.SetLogger(&MyLogger{})
}

type MyLogger struct{}

func (m *MyLogger) Log(sckey string, title string, code int32, err error) {}
```
