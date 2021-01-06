package serverchan

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	// serverchanApiUrl is the serverchan api's url, including sckey string.
	serverchanApiUrl = "https://sc.ftqq.com/%s.send"

	// contentType is the content type for serverchan api.
	contentType = "application/x-www-form-urlencoded"
)

// Client represents a serverchan client. Please visit http://sc.ftqq.com/3.version for details.
type Client struct{}

// NewClient creates a default Client.
func NewClient() *Client {
	return &Client{}
}

var (
	// ErrEmptyTitle is an error for empty title.
	ErrEmptyTitle = errors.New("serverchan: empty title")

	// ErrBadPushToken is an error for bad push token.
	ErrBadPushToken = errors.New("serverchan: bad push token")

	// ErrDuplicateMessage is an error for duplicate message.
	ErrDuplicateMessage = errors.New("serverchan: duplicate message")

	// ErrNotSuccess is an error for not success, used when non-json and errno-non-zero response.
	ErrNotSuccess = errors.New("serverchan: respond with not success")
)

// Send sends a message with title to serverchan using given sckey.
func (c *Client) Send(sckey, title, message string) error {
	sckey = strings.TrimSpace(sckey)
	title = strings.TrimSpace(title)
	message = strings.TrimSpace(message)
	if title == "" {
		return ErrEmptyTitle
	}
	if sckey == "" {
		return ErrBadPushToken
	}

	form := url.Values{}
	form.Add("text", title)
	form.Add("desp", message)

	req, err := http.NewRequest("POST", fmt.Sprintf(serverchanApiUrl, sckey), strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)

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

// CheckSckey sends a test message to serverchan and checks valid sckey. (through errno == 1024)
func (c *Client) CheckSckey(sckey string) (bool, error) {
	err := c.Send(sckey, "CheckSckey", "")

	if err == ErrBadPushToken {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// responseObject represents serverchan response object.
// Example responses:
// 	{"errno":0,"errmsg":"success","dataset":"done"}
// 	{"errno":1024,"errmsg":"\u4e0d\u8981\u91cd\u590d\u53d1\u9001\u540c\u6837\u7684\u5185\u5bb9"}
// 	{"errno":1024,"errmsg":"bad pushtoken"}
// 	<h2>系统消息</h2><p>消息标题不能为空</p>
type responseObject struct {
	Errno   int32  `json:"errno"`
	Errmsg  string `json:"errmsg"`
	Dataset string `json:"dataset"`
}

const (
	errnoSuccess       = 0
	errmsgBadPushToken = "bad pushtoken"
	errmsgDuplicate    = "不要重复发送同样的内容"
)
