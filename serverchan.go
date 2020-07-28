package serverchan

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	serverChanApi = "https://sc.ftqq.com/%s.send?text=%s&desp=%s"
)

type Client struct {
	logger Logger
}

// Serverchan logger, please override `Log` and `LogMode`.
type Logger interface {
	Log(sckey string, title string, err error)
	LogMode() bool
}

// Set serverchan's logger, use `DefaultLogger` or `NoLogger`.
func (c *Client) SetLogger(logger Logger) {
	c.logger = logger
}

// Send title and message to serverchan through `sc.ftqq.com`.
func (c *Client) Send(sckey string, title string, message string) (*http.Response, error) {
	title = url.QueryEscape(title)
	message = url.QueryEscape(message)
	sendUrl := fmt.Sprintf(serverChanApi, sckey, title, message)

	resp, err := http.Post(sendUrl, "application/x-www-form-urlencoded", strings.NewReader("name=cjb"))
	if err != nil {
		c.logger.Log(sckey, title, err)
		return resp, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err := fmt.Errorf("success to send but get %d response", resp.StatusCode)
		c.logger.Log(sckey, title, err)
		return resp, err
	}

	c.logger.Log(sckey, title, nil)
	return resp, nil
}

// A Serverchan default logger with log.Logger.
type DefaultLogger struct{}

func (d *DefaultLogger) Log(sckey string, title string, err error) {
	if err != nil {
		log.Printf("[Serverchan] failed to send message to %s: %v", sckey, err)
	} else {
		log.Printf("[Serverchan] <- %s | %s", sckey, title)
	}
}

func (d *DefaultLogger) LogMode() bool {
	return true
}

// A Serverchan logger that not to do anything.
type NoLogger struct{}

func (n *NoLogger) Log(sckey string, title string, err error) {}

func (n *NoLogger) LogMode() bool {
	return false
}
