package utils

// TruncateString - truncates a string to a given length and adds ... if the string exceeds the expected length
func TruncateString(str string, num int) string {
	truncated := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		truncated = str[0:num] + "..."
	}
	return truncated
}
