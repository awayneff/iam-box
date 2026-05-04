// app/helpers/cmd.go

package helpers

import "strconv"

func StrToInt(s string, defaultVal int) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}

	return num
}
