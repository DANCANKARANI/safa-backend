package services

import "regexp"

func ValidateEmail(email string) bool {
	// Basic RFC 5322 email regex (simplified)
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	if len(email) < 5 || len(email) > 254 {
		return false
	}

	return emailRegex.MatchString(email)
}
func ValidatePhoneNumber(phone string) bool {
	if len(phone) < 10 || len(phone) > 15 {
		return false
	}
	for _, char := range phone {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}
