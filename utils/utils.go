package utils

func MaskToken(token string) string {
	if len(token) <= 15 {
		return token
	}
	return token[:8] + "..." + token[len(token)-7:]
}
