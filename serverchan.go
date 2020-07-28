package serverchan

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

// Set serverchan's logger, use `DefaultLogger` or `NoLogger` or others.
func (c *Client) SetLogger(logger Logger) {
	c.logger = logger
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
	content := string(bs)
	if content == "" {
		return &ResponseObject{Errno: ErrnoDefault, Errmsg: errmsgDuplicate}, resp, nil
	}

	obj := &ResponseObject{}
	err = json.Unmarshal(bs, obj)
	if err != nil {
		return nil, resp, err
	}

	// check code
	if obj.Errno != ErrnoSuccess {
		if obj.Errmsg == errmsgDuplicate {
			obj.Errmsg = ErrmsgDuplicate
		}
		err := newResponseError(obj)
		c.logger.Log(sckey, title, obj.Errno, err)
		return obj, resp, err
	} else {
		c.logger.Log(sckey, title, obj.Errno, nil)
		return obj, resp, nil
	}
}

// Send a test message to serverchan and check user is existed. (through errno == 1024)
func (c *Client) CheckSckey(sckey string, testTitle string) (bool, error) {
	obj, _, err := c.Send(sckey, testTitle, "")
	if obj != nil {
		return obj.Errmsg != ErrmsgBadPushToken, nil
	}
	return false, err
}
