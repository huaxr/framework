// Author: huaxr
// Time: 2022-10-31 11:23
// Git: huaxr

package selector

import (
	"testing"
)

func TestSelector(t *testing.T) {
	x := NewSelector(RoundRobin, []string{"1", "2", "3"})

	go func() {
		for {
			tx := x.Select()
			t.Log(tx)
		}
	}()

	go func() {
		for {
			tx := x.Select()
			t.Log(tx)
		}
	}()

	select {}
}
