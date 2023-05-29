// Author: huaxr
// Time: 2022-10-22 14:42
// Git: huaxr

package httputil

import (
	"context"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
)

var cli = NewHttpClient(10*time.Second, 1, true)

func TestPost(t *testing.T) {
	testing.Init()
	res, err := cli.Post(context.Background(), "http://127.0.0.1:8888/grpcx_test/22", map[string]interface{}{"a": 111})
	if err != nil {
		t.Log("err", err)
		return
	}

	t.Log("success", string(res))
}

func get(t *testing.T) {
	//var cli = NewHttpClient(10*time.Second, 1, true)
	res, err := resty.New().R().Get("http://127.0.0.1:8888/test/11")
	if err != nil {
		t.Log("err", err)
		return
	}
	t.Log("success", res)
}

func TestGet(t *testing.T) {
	testing.Init()
	var i = 1000
	for i > 0 {
		i--
		go func() {
			get(t)
		}()
	}
	time.Sleep(1 * time.Second)
	get(t)
	select {}
}
