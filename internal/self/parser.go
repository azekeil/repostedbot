package self

import "strings"

func getCommand(s string) string {
	spl := strings.Split(s, " ")
	if len(spl) >= 2 {
		return spl[1]
	}
	return ""
}

func getArgs(s string) string {
	spl := strings.SplitAfterN(s, " ", 3)
	if len(spl) >= 3 {
		return spl[2]
	}
	return ""
}
