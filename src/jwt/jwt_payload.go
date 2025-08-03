package jwt

type JwtPayload struct {
	UserId         int64
	Username       string
	ExpirationTime int64
}
