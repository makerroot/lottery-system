package constants

// Validation rule constants
const (
	// Password requirements
	MinPasswordLength = 6
	MaxPasswordLength = 100

	// Username requirements
	MinUsernameLength = 3
	MaxUsernameLength = 50

	// Name requirements
	MinNameLength = 1
	MaxNameLength = 100

	// Phone number
	PhoneNumberLength = 11

	// Company code
	DefaultCompanyCode   = "DEFAULT"
	MaxCompanyCodeLength = 50
	MaxCompanyNameLength = 200

	// Prize levels
	MaxPrizeLevelNameLength = 50
	MaxPrizeDescLength      = 200

	// Pagination
	DefaultPageSize = 20
	MaxPageSize     = 100

	// Lottery
	DefaultDrawCount = 1
	MaxDrawCount     = 100
)

// Regular expression patterns for validation
const (
	// Phone number pattern (Chinese mobile)
	PhonePattern = `^1[3-9]\d{9}$`

	// Username pattern (alphanumeric and underscore)
	UsernamePattern = `^[a-zA-Z0-9_]+$`

	// Email pattern
	EmailPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Company code pattern
	CompanyCodePattern = `^[a-zA-Z0-9_-]+$`
)
