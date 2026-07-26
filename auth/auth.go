package auth

type Authenticator interface {
	Authenticate(name string, token []byte) error
}
