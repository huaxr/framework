// Author: huaxr
// Time: 2022-10-26 10:12
// Git: huaxr

package define

import (
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestValidateBasicPath(t *testing.T) {
	convey.Convey("TestValidateBasicPath", t, func() {
		var p1 = gomonkey.ApplyFunc(ValidateBasicPath, func(srv string) error {
			return nil
		})
		defer p1.Reset()
		err := ValidateBasicPath("a.a.")
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestPath(t *testing.T) {
	p := ValidateBasicPath("/aa___aaaa..a.aa1")
	t.Log(p)

	p = ValidateBasicPath("/a.a.")
	t.Log(p)

	p = ValidateBasicPath("/a.a.aa1.")
	t.Log(p)

	p = ValidateBasicPath("/..")
	t.Log(p)

	p = ValidateBasicPath("/x..x")
	t.Log(p)

	p = ValidateBasicPath("/a.a.aa1")
	t.Log(p)
}
