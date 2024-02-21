package interfaces

import (
	"context"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/Zeta-Manu/Backend/internal/domain/entity"
)

type IIdentityProvider interface {
	Register(ctx context.Context, userRegistration entity.UserRegistration) (string, error)
	Login(ctx context.Context, userLogin entity.UserLogin) (string, error)
	GetUser(ctx context.Context, email string) (*cognitoidentityprovider.AdminGetUserOutput, error)
	ConfirmRegistration(ctx context.Context, userRegistrationConfirm entity.UserRegistrationConfirm) error
	ResendConfirmationCode(ctx context.Context, email entity.Email) (string, error)
	ForgotPassword(ctx context.Context, email entity.Email) (string, error)
	ConfirmForgotPassword(ctx context.Context, userResetPassword entity.UserResetPassword) (*cognitoidentityprovider.ConfirmForgotPasswordOutput, error)
	ChangePassword(ctx context.Context, accessToken string, userChangePassword entity.UserChangePassword) error
}
