package string

import (
	"strings"
)

func RemoveSpecialCharacters(str string) string {
	str = strings.ReplaceAll(str, ".", "")
	str = strings.ReplaceAll(str, "-", "")
	str = strings.ReplaceAll(str, "/", "")

	return str
}

func ExtractTokenFromString(str string) string {
	propList := strings.Split(str, "\n")

	for _, value := range propList {
		key := strings.Split(value, ":")
		if key[0] == "Authorization" {
			key[1] = strings.ReplaceAll(key[1], "Bearer ", "")
			key[1] = strings.ReplaceAll(key[1], " ", "")
			if key[1] != "" {
				return key[1]
			}
		}
	}

	return ""
}

func FormatGenericMessagesString(strGeneric, strTarget string) string {
	return strings.ReplaceAll(strGeneric, "%_%", strTarget)
}
