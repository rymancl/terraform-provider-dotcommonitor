package dotcommonitor

import "strings"

// StateToLower ... converts string to lowercase to normalize it for state
func StateToLower(v interface{}) string {
	s, ok := v.(string)

	if !ok {
		return ""
	}

	return strings.ToLower(s)
}
