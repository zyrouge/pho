package utils

import (
	"fmt"
	"time"
)

func TimeNowSeconds() int64 {
	return time.Now().Unix()
}

func HumanizeSeconds(n int64) string {
	s := n % 60
	n /= 60
	m := n % 60
	h := n / 60
	if h == 0 && m == 0 {
		return fmt.Sprintf("%ds", s)
	}
	if h == 0 {
		return fmt.Sprintf("%dm:%ds", m, s)
	}
	return fmt.Sprintf("%dh:%dm:%ds", h, m, s)
}
