package shared

// Business error codes
const (
	// 1000-1999: Validation errors
	ErrCodeValidation  = 1000
	ErrCodeBadRequest  = 1001
	ErrCodeInvalidSlug = 1002

	// 2000-2999: Authentication errors
	ErrCodeUnauthorized     = 2000
	ErrCodeInvalidToken     = 2001
	ErrCodeTokenExpired     = 2002
	ErrCodeInvalidCredentials = 2003
	ErrCodeForbidden        = 2004

	// 4000-4999: Not found errors
	ErrCodeNotFound      = 4000
	ErrCodePostNotFound  = 4001
	ErrCodeUserNotFound  = 4002
	ErrCodeTagNotFound   = 4003
	ErrCodeProjectNotFound = 4004

	// 5000-5999: Server errors
	ErrCodeInternal     = 5000
	ErrCodeDatabase     = 5001
	ErrCodeExternal     = 5002
)
