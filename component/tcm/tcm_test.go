package tcm

import (
	"context"
	"fmt"
	"testing"
)

func TestGetConfig(t *testing.T) {
	testing.Init()
	InitTcmInstance(context.Background())
	test := GetTextConfig("test")
	fmt.Println("test:", test)
	test2 := GetTextConfig("")
	fmt.Println("test2:", test2)
	select {}
}
