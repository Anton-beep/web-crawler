package utils

import (
	"regexp"
	"strings"
)

// IsCorrectLink checks if the link is a valid URL
func IsCorrectLink(link string) bool {
	if strings.Count(link, "http://")+strings.Count(link, "https://") == 1 {
		regex := `^((http|https):\/\/)(localhost|[a-zA-Z0-9.-]+)(:\d+)?(\/[a-zA-Z0-9@:%._\+~#?&//=]*)?$`
		match, _ := regexp.MatchString(regex, link)
		return match
	}
	return false
}
