package validation

import (
	"fmt"
	"regexp"
	"strings"
)

func ValidateEmail(email string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return regex.MatchString(email)
}

func ValidatePhone(phone string) bool {
	regex := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	return regex.MatchString(phone)
}

func SanitizeInput(input string) string {
	input = strings.TrimSpace(input)
	re := regexp.MustCompile(`\s+`)
	input = re.ReplaceAllString(input, " ")
	return input
}

func ValidateUpdateProfileRequest(name, email string) error {
	if name != "" && (len(name) < 2 || len(name) > 100) {
		return fmt.Errorf("name must be between 2 and 100 characters")
	}

	if email != "" && !ValidateEmail(email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}
