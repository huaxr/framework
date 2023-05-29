package notifyutil

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/url"
	"time"

	"github.com/huaxr/framework/logx"
	"github.com/huaxr/framework/pkg/httputil"
	"github.com/spf13/cast"
)

// https://doc-openapi.saash.vdyoo.com/doc#/robot/develop?id=%e8%87%aa%e5%ae%9a%e4%b9%89%e6%9c%ba%e5%99%a8%e4%ba%ba
const HOST = "https://xxx.com/robot/send"

type RobotSendRequest struct {
	AccessToken string `json:"access_token"`
	SecretKey   string `json:"secret_key"`
	Message     string `json:"message"`
}

type RobotSendResponse struct {
	YachMid interface{} `json:"yachMid"`
}

// 自定义机器人
func RobotSend(ctx context.Context, req RobotSendRequest) (resp RobotSendResponse, err error) {
	timestamp := time.Now().UnixNano() / 1e6
	tmpData := cast.ToString(timestamp) + "\n" + req.SecretKey
	hmaData := computeHmacSha256(tmpData, req.SecretKey)
	sign := url.QueryEscape(hmaData)
	reqPath := HOST + "?access_token=" + req.AccessToken + "&timestamp=" + cast.ToString(timestamp) + "&sign=" + sign

	var reqMsg map[string]interface{}
	_ = json.Unmarshal([]byte(req.Message), &reqMsg)

	cli := httputil.NewHttpClient(5*time.Second, 3, false)

	res, err := cli.Post(context.Background(), reqPath, reqMsg)
	if err != nil {
		logx.L().Debugf("post err:%v", err)
		return
	}

	err = json.Unmarshal(res, &resp)
	if err != nil {
		return
	}
	return
}

func computeHmacSha256(data string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	sha := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(sha)
}
