// Author: huaxinrui@tal.com
// Time:   2022/1/6 上午11:23
// Git:    huaxr

package toolutil

import (
	"errors"
	"log"
)

func CheckStr(str string) {
	for i, j := range str {
		r := rune(j)
		log.Println(i, j, r)
	}
}

func AZaz09_(str string) error {
	for _, i := range str {
		r := rune(i)
		if (r >= 65 && r <= 90) || (r >= 97 && r <= 122) || (r >= 48 && r <= 57) || r == 95 {
			continue
		} else {
			return errors.New("only support A-Z,a-z,0-9,_")
		}
	}
	return nil
}

func Str2Int(str string) int {
	var tmp int32
	for _, i := range str {
		tmp += i
	}
	return int(tmp)
}
