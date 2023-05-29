package notifyutil

import (
	"context"
	"fmt"
	"testing"
)

func Test_RobotSend(t *testing.T) {
	var req RobotSendRequest
	req.AccessToken = "ZmdHQ3pBeTdqYnlNOXhlZVl4MjJtTE1PeGU0eE84eTBpbWhsL0tLa1Y3dnF6azlyUzRKSmMwRHUyanlTTjFMQg"
	req.SecretKey = "SEC9d1266d83c7a45ea0074bbebdb4398d0"
	msg := `{"msgtype":"text","text":{"content":"通知测试"}}`
	req.Message = msg
	resp, err := RobotSend(context.Background(), req)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp)
}
