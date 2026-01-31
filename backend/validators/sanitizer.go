package validators

import (
	"html"
	"regexp"
	"strings"
	"unicode"
)

// SanitizeString sanitizes a string input by removing potentially dangerous characters
func SanitizeString(input string) string {
	if input == "" {
		return ""
	}

	// Trim whitespace
	input = strings.TrimSpace(input)

	// Remove HTML tags
	input = html.EscapeString(input)

	// Remove control characters (except newline, tab, carriage return)
	input = removeControlCharacters(input)

	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")

	return input
}

// SanitizeUsername sanitizes a username
func SanitizeUsername(username string) string {
	username = strings.TrimSpace(username)

	// Remove any characters that aren't alphanumeric or underscore
	reg := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	username = reg.ReplaceAllString(username, "")

	return username
}

// SanitizeName sanitizes a person's name
func SanitizeName(name string) string {
	name = strings.TrimSpace(name)

	// Remove control characters
	name = removeControlCharacters(name)

	// Remove null bytes
	name = strings.ReplaceAll(name, "\x00", "")

	return name
}

// SanitizePhone sanitizes a phone number
func SanitizePhone(phone string) string {
	phone = strings.TrimSpace(phone)

	// Remove all non-digit characters
	reg := regexp.MustCompile(`[^\d]`)
	phone = reg.ReplaceAllString(phone, "")

	return phone
}

// SanitizeEmail sanitizes an email address
func SanitizeEmail(email string) string {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)

	// Remove null bytes
	email = strings.ReplaceAll(email, "\x00", "")

	return email
}

// SanitizeSQL attempts to detect SQL injection patterns
func SanitizeSQL(input string) (string, bool) {
	if input == "" {
		return input, true
	}

	// Common SQL injection patterns
	patterns := []string{
		`(?i)(\bunion\b.*\bselect\b)`,
		`(?i)(\bselect\b.*\bfrom\b)`,
		`(?i)(\binsert\b.*\binto\b)`,
		`(?i)(\bdelete\b.*\bfrom\b)`,
		`(?i)(\bdrop\b.*\btable\b)`,
		`(?i)(\bexec\b|\bexecute\b)`,
		`(?i)(\bscript\b.*:)`,
		`(;|\b(or\b|and\b).*(=|==))`,
		`(?i)(\bdeclare\b.*@\b)`,
		`(?i)(\bcast\b.*\bexec\b)`,
		`['";]`,
	}

	inputLower := strings.ToLower(input)
	for _, pattern := range patterns {
		matched, _ := regexp.MatchString(pattern, inputLower)
		if matched {
			return "", false // Potential SQL injection detected
		}
	}

	return input, true
}

// removeControlCharacters removes control characters except \n, \r, \t
func removeControlCharacters(s string) string {
	var result []rune
	for _, r := range s {
		if r == '\n' || r == '\r' || r == '\t' {
			result = append(result, r)
		} else if !unicode.IsControl(r) {
			result = append(result, r)
		}
	}
	return string(result)
}

// ValidateXSS attempts to detect XSS patterns
func ValidateXSS(input string) bool {
	if input == "" {
		return true
	}

	// Common XSS patterns
	patterns := []string{
		`<script[^>]*>.*?</script>`,
		`javascript:`,
		`on\w+\s*=`, // Events like onclick=
		`<iframe`,
		`<embed`,
		`<object`,
		`<link`,
		`<meta`,
		`<style`,
		`<applet`,
		`<body`,
	}

	inputLower := strings.ToLower(input)
	for _, pattern := range patterns {
		matched, _ := regexp.MatchString(pattern, inputLower)
		if matched {
			return true // Potential XSS detected
		}
	}

	return false
}
