# go-serverchan

+ Unofficial serverchan client library in Golang.

### ServerChan

+ Please visit [Server酱](http://sc.ftqq.com/3.version) for details.

> Attention:
>
> 不要在text参数中传递引号、点、花括号等字符。因为微信的接口不支持一系列的特殊字符，但没有详细列表，所以我只简单的过滤掉了一些。
> 如果需要发送特殊字符，请放到 desp字段中。
>
> 每人每天发送上限500条，相同内容5分钟内不能重复发送，不同内容一分钟只能发送30条。主要是防止程序出错的情况。
> 注意，因为之前频繁的死循环程序导致费用飙升，现在每天调用接口超过1000次的用户将被系统自动拉黑。

### Usage

+ See [example_test.go](./example_test.go) for more demos.

```go
package main

import (
	"context"
	"errors"
	"time"

	serverchan "github.com/Aoi-hosizora/go-serverchan"
)

func main() {
	client := serverchan.NewClient()

	err := client.Send("xxx", "test title", "test message")
	if err != nil {
		// ...
	}

	ok, err := client.CheckSckey("xxx", "check sckey")
	if err != nil {
		// ...
	} else if !ok {
		// ...
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = client.SendWithContext(ctx, "xxx", "test title", "test message")
	if errors.Is(err, context.DeadlineExceeded) {
		// ...
	}
}
```
