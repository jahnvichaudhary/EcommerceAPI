package product

import "strconv"

func Float64ToStr(num float64) string {
	return strconv.FormatFloat(num, 8, 2, 8)
}
