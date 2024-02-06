package utils

import "strings"

func RemoveSpecialChars2(s string) string {
	r := strings.NewReplacer("\t", "", "\n", " ", "\r", "", "\x00", "")
	return r.Replace(s)
}

var Layout = "2006-01-02 15:04:05"
