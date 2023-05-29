// Author: huaxr
// Time: 2022-11-04 13:38
// Git: huaxr

package version

import (
	"fmt"

	"github.com/spf13/cast"
)

const Version = 1399

// GetVersionStr return v1.x.xx
func GetVersionStr() string {
	v := cast.ToString(Version)
	return fmt.Sprintf("v%c.%c.%c%c", v[0], v[1], v[2], v[3])
}
