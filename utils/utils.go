package utils

func MangleToken(token string) string {
	if len(token) <= 20 {
		return token
	}
	return token[:10] + "..." + token[len(token)-10:]
}
