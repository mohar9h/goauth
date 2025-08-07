package auth

// Result TokenResult is the result of a successful token creation.
type Result struct {
	PlainText string // what the client receives
	TokenID   string // internal hashed ID for storage
}
