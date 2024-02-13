package repository

import (
	"context"

	"github.com/Zeta-Manu/Backend/internal/adapters/identityprovider"
)

type CognitoRepository struct {
	provider identityprovider.IdentityProvider
}

func NewCognitoRepository(provider identityprovider.IdentityProvider) *CognitoRepository {
	return &CognitoRepository{provider: provider}
}

func (r *CognitoRepository) RegisterUser(ctx context.Context, name, email, password string) (string, error) {
	return r.provider.Register(ctx, name, email, password)
}

func (r *CognitoRepository) LoginUser(ctx context.Context, email, password string) (string, error) {
	return r.provider.Login(ctx, email, password)
}
