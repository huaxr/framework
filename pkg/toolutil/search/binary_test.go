// Author: XinRui Hua
// Time:   2023/01/03 18:12
// Git:    huaxr

package search

import (
	"fmt"
	"testing"
)

func TestBinFind(t *testing.T) {
	arr := make([]int, 1024*1024, 1024*1024)
	for i := 0; i < 1024*1024; i++ {
		arr[i] = i + 1
	}
	id := binSearch(arr, 1024)
	if id != -1 {
		fmt.Println(id, arr[id])
	} else {
		fmt.Println("no data")
	}
}
