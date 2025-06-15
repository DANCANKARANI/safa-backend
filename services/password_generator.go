package services

import (
	"crypto/rand"
	"fmt"
	"strconv"
	"time"
)

// GenerateFormattedPassword returns a password like "1234.2025"
func GenerateFormattedPassword() (string, error) {
	const digitCount = 4
	var digits string

	// Generate 4 random digits
	for len(digits) < digitCount {
		b := make([]byte, 1)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}
		n := int(b[0]) % 10
		digits += strconv.Itoa(n)
	}

	// Get current year
	year := time.Now().Year()

	return fmt.Sprintf("%s.%d", digits, year), nil
}
