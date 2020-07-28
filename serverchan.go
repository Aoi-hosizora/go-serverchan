package serverchan

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func NewClient() *Client {
	return &Client{logger: DefaultLogger(LMMask)}
}

// Serverchan logger, please override `Log`.
type Logger interface {
	Log(sckey string, title string, code int32, err error)
}

// Set serverchan's logger, use `DefaultLogger` or `NoLogger`.
func (c *Client) SetLogger(logger Logger) {
	c.logger = logger
}

// Serverchan's response model
type ResponseObject struct {
	// some code: 0, 1024
	Errno int32 `json:"errno"`

	// some message: success, bad pushtoken
	Errmsg string `json:"errmsg"`

	// demo dataset: done
	Dataset string `json:"dataset"`
}

// Send a test message to serverchan and check user is existed. (through errno == 1024)
func (c *Client) CheckExist(sckey string, testTitle string) (bool, error) {
	obj, _, err := c.Send(sckey, testTitle, "")
	if obj != nil {
		return obj.Errno != 1024, err
	}
	return false, err
}

// Send title and message to serverchan through `sc.ftqq.com`.
func (c *Client) Send(sckey string, title string, message string) (*ResponseObject, *http.Response, error) {
	title = url.QueryEscape(title)
	message = url.QueryEscape(message)
	sendUrl := fmt.Sprintf(serverChanApi, sckey, title, message)

	// check title
	if title == "" {
		err := fmt.Errorf("title could not be empty")
		c.logger.Log(sckey, title, -1, err)
		return nil, nil, err
	}

	// post http
	resp, err := http.Post(sendUrl, "application/x-www-form-urlencoded", strings.NewReader("name=cjb"))
	if err != nil {
		c.logger.Log(sckey, title, -1, err)
		return nil, nil, err
	}
	body := resp.Body
	defer body.Close()

	// response content
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, err
	}
	obj := &ResponseObject{}
	err = json.Unmarshal(bs, obj)
	if err != nil {
		return nil, resp, err
	}

	// check code
	if obj.Errno != 0 {
		err := fmt.Errorf("%d: %s", obj.Errno, obj.Errmsg)
		c.logger.Log(sckey, title, obj.Errno, err)
		return obj, resp, err
	} else {
		c.logger.Log(sckey, title, obj.Errno, nil)
		return obj, resp, nil
	}
}

// Serverchan log mode, use this to control log behavior.
type LogMode uint8

const (
	// Not to log.
	LMNone LogMode = iota

	// Log error only.
	LMErr

	// Log all, but with masked sckey and title.
	LMMask

	// Log all, include full sckey.
	LMAll
)

type defaultLogger struct {
	mode LogMode
}

// A serverchan default logger with log.Logger.
func DefaultLogger(mode LogMode) *defaultLogger {
	return &defaultLogger{mode: mode}
}

func (d *defaultLogger) Log(sckey string, title string, code int32, err error) {
	mode := d.mode
	if mode <= LMNone {
		return
	}

	if mode >= LMErr {
		// Err, Mask, All
		if err != nil {
			log.Printf("[Serverchan] failed to send message to %s: %v", Mask(sckey), err)
			return
		}
	}

	if mode >= LMMask {
		// Mask, All
		if mode == LMMask {
			sckey = Mask(sckey)
			title = Mask(title)
		}

		if code == 0 {
			log.Printf("[Serverchan] <- %3d | %s | %s", 0, sckey, title)
		}
	}
}

type noLogger struct{}

// A serverchan logger (LMNone) that not to do anything.
func NoLogger() *noLogger {
	return &noLogger{}
}

// noinspection GoUnusedParameter
func (n *noLogger) Log(sckey string, title string, code int32, err error) {}

// Mask sckey as `*`.
func Mask(tok string) string {
	switch len(tok) {
	case 0:
		return ""
	case 1:
		return "*"
	case 2:
		return "*" + tok[1:]
	case 3:
		return "**" + tok[2:3]
	case 4:
		return tok[0:1] + "**" + tok[3:4]
	case 5:
		return tok[0:1] + "***" + tok[4:5]
	default:
		return tok[0:2] + strings.Repeat("*", len(tok)-4) + tok[len(tok)-2:] // <<< Default
	}
}
