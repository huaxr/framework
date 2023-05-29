// Author: huaxr
// Time: 2022-11-07 13:14
// Git: huaxr

package confutil

import (
	"fmt"

	"github.com/huaxr/framework/internal/define"
)

type PSM string

func (p PSM) Validate() error {
	err := define.ValidateBasicPath(fmt.Sprintf("/%s", p))
	if err != nil {
		return err
	}
	return nil
}

func (p PSM) String() string {
	return string(p)
}

func (p PSM) Equals(psm string) bool {
	return p.String() == psm
}
