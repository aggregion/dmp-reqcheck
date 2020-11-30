package utils

// Coalesce return first non empty argument
func Coalesce(strArgs ...string) string {
	for _, str := range strArgs {
		if len(str) > 0 {
			return str
		}
	}
	return ""
}
