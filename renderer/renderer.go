package renderer

import "strings"

func BooleanIcon(flag bool) string {
	if flag {
		return "✔"
	}

	return "✖"
}

func JoinStrings(allStrings []string) string {
	return strings.Join(allStrings, ", ")
}

func LimitString(fullString string, maxLength int) string {
	if len(fullString) > maxLength {
		return fullString[0:maxLength] + "..."
	}
	return fullString
}
