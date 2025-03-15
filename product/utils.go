package product

import "strconv"

func FloatToString(num float64) string {
	return strconv.FormatFloat(num, 8, 2, 8)
}

func StringToFloat(num string) float64 {
	res, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return 0
	}
	return res
}
