package interfaces

import (
	"context"

	"github.com/Zeta-Manu/Backend/internal/domain/entity"
)

type IIdentityProvider interface {
	Register(ctx context.Context, userRegistration entity.UserRegistration) (string, error)
	Login(ctx context.Context, userLogin entity.UserLogin) (string, error)
}
