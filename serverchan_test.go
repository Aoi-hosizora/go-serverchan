package serverchan

import (
	"log"
	"testing"
)

func TestServerChan(t *testing.T) {
	client := NewClient()
	log.Println(client.CheckExist("xxx", "test"))

	client.SetLogger(NoLogger())
	_, _, _ = client.Send("xxx", "title", "message")

	client.SetLogger(DefaultLogger(LMErr))
	_, _, _ = client.Send("xxx", "title", "message")

	client.SetLogger(DefaultLogger(LMMask))
	_, _, _ = client.Send("xxx", "title", "message")

	client.SetLogger(DefaultLogger(LMAll))
	_, _, _ = client.Send("xxx", "title", "message")
}
