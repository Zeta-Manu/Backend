package identityprovider

import (
	"context"
)

type IdentityProvider interface {
	Register(ctx context.Context, name, email, password string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
}
