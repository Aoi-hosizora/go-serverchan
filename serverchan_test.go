package serverchan

import (
	"log"
	"testing"
)

func TestServerChan(t *testing.T) {
	client := NewClient()
	log.Println(client.CheckSckey("xxx", "test"))

	client.SetLogger(NoLogger())
	obj, _, err := client.Send("xxx", "title", "message")
	log.Println(obj, IsResponseError(err))

	client.SetLogger(DefaultLogger(LMErr))
	_, _, _ = client.Send("xxx", "title", "message")

	client.SetLogger(DefaultLogger(LMMask))
	_, _, _ = client.Send("xxx", "title", "message")

	client.SetLogger(DefaultLogger(LMAll))
	_, _, _ = client.Send("xxx", "title", "message")
}
