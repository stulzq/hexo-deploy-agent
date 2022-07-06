package util

import (
	"os"
	"strings"
)

func IsDebug() bool {
	val := os.Getenv("HEXO_DEBUG")

	if strings.ToLower(val) == "false" {
		return false
	}

	return true
}
