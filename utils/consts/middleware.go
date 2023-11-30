package consts

const (
	JWTContextKey     = "idTokenCustomClaims"
	JWTUserContextKey = "user"

	ErrBadRequestJWT   = "Missing or malformed JWT"
	ErrUnauthorizedJWT = "Invalid or expired JWT"
)
