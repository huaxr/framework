package define

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var reg = regexp.MustCompile(`^/[\w]+\.{1}[\w]+\.{1}[\w]+$`)

func ValidateBasicPath(srv string) error {
	if !strings.HasPrefix(srv, "/") {
		return errors.New("srvPath must start with `/`")
	}

	if len(srv) < 8 || len(srv) > 50 {
		return errors.New("grpc's path must define a name at least 8 character and less than 50")
	}

	res := reg.FindAllString(srv, -1)
	if len(res) == 0 {
		return errors.New(fmt.Sprintf("invalide format:%v, you should define"+
			" psm by format of `P.S.M` (产品.服务.模块) rule, for instance: `montage.framework.test`", srv))
	}

	return nil
}

func ValidatePort(port int) error {
	if port < 5000 || port >= 10000 {
		return errors.New("should 5000<=port<10000")
	}
	return nil
}
