package envauther

import (
	"crypto/subtle"
	"os"

	"github.com/ajsqr/wombat/auth"
)

var _ auth.Authenticator = &EnvAuthenticator{}

type EnvAuthenticator struct {
}

func NewEnvAuthenticator() *EnvAuthenticator {
	return &EnvAuthenticator{}
}

func (e *EnvAuthenticator) Authenticate(name string, token []byte) error {
	realToken, ok := os.LookupEnv(name)
	if !ok {
		return auth.ErrTokenNotSet
	}

	if subtle.ConstantTimeCompare([]byte(realToken), token) == 1 {
		return nil
	}

	return auth.ErrAccessDenied
}
