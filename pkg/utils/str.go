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

// IsIntersectStrs .
func IsIntersectStrs(strArgs1 []string, strArgs2 []string) bool {
	for _, s1 := range strArgs1 {
		for _, s2 := range strArgs2 {
			if s1 == s2 {
				return true
			}
		}
	}

	return false
}
