package middleware

// ValidateToken takes an access token (JWT) and validates it,
// returning true if it's valid and false otherwise.
// TODO: Echo's own JWT middleware validator is probably enough.
// After implementing user login, decide this
func ValidateToken(accTkn string) (bool, error) {
	return false, nil
}
