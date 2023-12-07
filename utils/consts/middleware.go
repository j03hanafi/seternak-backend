package consts

const (
	JWTContextKey             = "idTokenCustomClaims"
	JWTUserContextKey         = "user"
	JWTRefreshTokenContextKey = "refreshToken"

	ErrBadRequestJWT   = "Missing or malformed JWT"
	ErrUnauthorizedJWT = "Invalid or expired JWT"
)
