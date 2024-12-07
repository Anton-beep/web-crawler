package utils

import (
	"regexp"
	"strings"
)

func IsCorrectLink(link string) bool {
	if strings.Count(link, "http://")+strings.Count(link, "https://") == 1 {
		regex := `^((http|https):\/\/)(localhost|[a-zA-Z0-9.-]+)(:\d+)?(\/[a-zA-Z0-9@:%._\+~#?&//=]*)?$`
		match, _ := regexp.MatchString(regex, link)
		return match
	}
	return false
}
