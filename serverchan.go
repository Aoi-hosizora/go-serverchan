package serverchan

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	serverchanApiUrl = "https://sc.ftqq.com/%s.send"       // serverchan api's url
	contentType      = "application/x-www-form-urlencoded" // content type for serverchan api
)

var (
	ErrEmptyTitle       = errors.New("serverchan: empty title")              // Empty title.
	ErrBadPushToken     = errors.New("serverchan: bad push token")           // Bad push token.
	ErrDuplicateMessage = errors.New("serverchan: duplicate message")        // Duplicate message.
	ErrNotSuccess       = errors.New("serverchan: respond with not success") // Not success, used when respond non-json or non-zero errno.
)

// Client represents a serverchan client. Please visit http://sc.ftqq.com/3.version for details.
type Client struct{}

// NewClient creates a default Client.
func NewClient() *Client {
	return &Client{}
}

// Send sends a message to serverchan using given sckey.
func (c *Client) Send(sckey, title, message string) error {
	sckey, form, err := checkParameters(sckey, title, message)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf(serverchanApiUrl, sckey), strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)

	return doRequest(req)
}

// SendWithContext sends a message to serverchan using given sckey (with context).
func (c *Client) SendWithContext(ctx context.Context, sckey, title, message string) error {
	sckey, form, err := checkParameters(sckey, title, message)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf(serverchanApiUrl, sckey), strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)

	return doRequest(req)
}

// CheckSckey sends a test message to serverchan to check if the sckey is valid.
func (c *Client) CheckSckey(sckey, title string) (bool, error) {
	err := c.Send(sckey, title, "")
	if err == ErrBadPushToken {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// CheckSckeyWithContext sends a test message to serverchan to check if the sckey is valid (with context).
func (c *Client) CheckSckeyWithContext(ctx context.Context, sckey, title string) (bool, error) {
	err := c.SendWithContext(ctx, sckey, title, "")
	if err == ErrBadPushToken {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// checkParameters checks parameters and returns new sckey and url.Values for http.Request.
func checkParameters(sckey, title, message string) (string, url.Values, error) {
	sckey = strings.TrimSpace(sckey)
	title = strings.TrimSpace(title)
	message = strings.TrimSpace(message)
	if sckey == "" {
		return "", nil, ErrBadPushToken
	}
	if title == "" {
		return sckey, nil, ErrEmptyTitle
	}

	form := url.Values{}
	form.Add("text", title)
	form.Add("desp", message)
	return sckey, form, nil
}

// doRequests sends the given request and parse the response.
func doRequest(req *http.Request) error {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body := resp.Body
	defer body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	obj := &responseObject{}
	err = json.Unmarshal(bs, obj)
	if err != nil {
		return ErrNotSuccess
	}
	if obj.Errmsg == errmsgBadPushToken {
		return ErrBadPushToken
	}
	if obj.Errmsg == errmsgDuplicate {
		return ErrDuplicateMessage
	}
	if obj.Errno != errnoSuccess {
		return ErrNotSuccess
	}

	return nil
}

// responseObject represents serverchan response object.
//
// Examples:
// 	{"errno":0,"errmsg":"success","dataset":"done"}
// 	{"errno":1024,"errmsg":"\u4e0d\u8981\u91cd\u590d\u53d1\u9001\u540c\u6837\u7684\u5185\u5bb9"}
// 	{"errno":1024,"errmsg":"bad pushtoken"}
// 	<h2>系统消息</h2><p>消息标题不能为空</p>
type responseObject struct {
	Errno   int32  `json:"errno"`
	Errmsg  string `json:"errmsg"`
	Dataset string `json:"dataset"`
}

// Constants related to responseObject.
const (
	errnoSuccess       = 0               // success errno
	errmsgBadPushToken = "bad pushtoken" // bad push token errmsg
	errmsgDuplicate    = "不要重复发送同样的内容"   // duplicate message errmsg
)
