// app/api/request.go

package api

import (
	"net/http"
	"strconv"
)

func ParseIntQuery(r *http.Request, key string, defaultVal int) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultVal
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	
	return intVal
}
