package jwt

type JwtPayload struct {
	Username       string
	ExpirationTime int64
}
