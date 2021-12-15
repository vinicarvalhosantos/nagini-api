package string

import "strings"

func RemoveSpecialCharacters(str string) string {
	str = strings.ReplaceAll(str, ".", "")
	str = strings.ReplaceAll(str, "-", "")
	str = strings.ReplaceAll(str, "/", "")
	str = strings.ReplaceAll(str, "(", "")
	str = strings.ReplaceAll(str, ")", "")

	return str
}

func FormatGenericMessagesString(strGeneric, strTarget string) string {
	return strings.ReplaceAll(strGeneric, "%_%", strTarget)
}
