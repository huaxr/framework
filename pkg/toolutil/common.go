// Author: huaxinrui@tal.com
// Time:   2021/12/31 上午10:51
// Git:    huaxr

package toolutil

import (
	"io/ioutil"
	"os/exec"
	"reflect"
	"time"
)

func GetRecentDayTimeStr(recent int) string {
	t := time.Now()
	nTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	yesTime := nTime.AddDate(0, 0, -recent)
	lastWeekDay := yesTime.Format("2006-01-02 15:04:05")
	return lastWeekDay
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func System(command string) string {
	cmd := exec.Command("/bin/bash", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {

	}
	defer stdout.Close()
	if err := cmd.Start(); err != nil {

	}
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {

	}
	return string(opBytes)
}
