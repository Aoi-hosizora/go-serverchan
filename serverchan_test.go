package serverchan

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestServerChan(t *testing.T) {
	form := url.Values{}
	form.Add("text", " ")
	form.Add("desp", "测试测试测试测试测试")
	req, _ := http.NewRequest("POST", "https://sc.ftqq.com/测.send", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	body := resp.Body
	defer body.Close()
	bs, err := ioutil.ReadAll(body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(resp.StatusCode, string(bs))
}
