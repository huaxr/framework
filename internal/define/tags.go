// Author: huaxr
// Time: 2022-11-16 11:25
// Git: huaxr

package define

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cast"
)

func Nano() string {
	return cast.ToString(time.Now().UnixNano())
}

func Uid() string {
	return uuid.NewV4().String()
}

func Chrome() string {
	return "chrome"
}

func Unknown() string {
	return "unknown"
}
