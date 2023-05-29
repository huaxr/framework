// Author: XinRui Hua
// Time:   2023/01/03 18:11
// Git:    huaxr

package search

import "fmt"

func binSearch(arr []int, findData int) int {
	low := 0
	high := len(arr) - 1
	for low <= high {
		mid := (low + high) / 2
		fmt.Println(mid)
		if arr[mid] > findData {
			high = mid - 1
		} else if arr[mid] < findData {
			low = mid + 1
		} else {
			return mid
		}
	}
	return -1
}
