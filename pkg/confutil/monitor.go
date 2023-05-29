// Author: huaxr
// Time: 2022-12-05 11:22
// Git: huaxr

package confutil

type Monitor struct {
	Url      string `yaml:"url"`
	Interval int    `yaml:"interval"` // s
}
