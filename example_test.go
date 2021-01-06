package serverchan_test

import (
	"context"
	"errors"
	"github.com/Aoi-hosizora/go-serverchan"
	"log"
	"time"
)

const sckey = "xxx"

func ExampleNewClient() {
	client := serverchan.NewClient()

	_ = client.Send(sckey, "test title", "test message")
	_, _ = client.CheckSckey("xxx", "check sckey")
}

func ExampleClient_Send() {
	client := serverchan.NewClient()

	err := client.Send(sckey, "test title", "# test message\n## sub title\n+ item1\n+ [item2](https://www.google.co.jp)")
	if err != nil {
		log.Println("Failed to send message:", err)
	} else {
		log.Println("Success to send message")
	}
}

func ExampleClient_SendWithContext() {
	client := serverchan.NewClient()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := client.SendWithContext(ctx, sckey, "test title", "test message")
	if errors.Is(err, context.DeadlineExceeded) {
		log.Println("Failed to send message: timeout")
	} else if err != nil {
		log.Println("Failed to send message:", err)
	} else {
		log.Println("Success to send message")
	}
}

func ExampleClient_CheckSckey() {
	client := serverchan.NewClient()

	ok, err := client.CheckSckey(sckey, "message for checking sckey")
	if err != nil {
		log.Println("Failed to send message:", err)
	} else if !ok {
		log.Println("Invalid sckey")
	} else {
		log.Println("Valid sckey")
	}
}

func ExampleClient_CheckSckeyWithContext() {
	client := serverchan.NewClient()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	ok, err := client.CheckSckeyWithContext(ctx, sckey, "message for checking sckey")
	if errors.Is(err, context.DeadlineExceeded) {
		log.Println("Failed to send message: timeout")
	} else if err != nil {
		log.Println("Failed to send message:", err)
	} else if !ok {
		log.Println("Invalid sckey")
	} else {
		log.Println("Valid sckey")
	}
}
