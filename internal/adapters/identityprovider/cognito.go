package identityprovider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/Zeta-Manu/Backend/internal/adapters/interfaces"
	"github.com/Zeta-Manu/Backend/internal/domain/entity"
)

var _ interfaces.IIdentityProvider = (*CognitoAdapter)(nil)

type CognitoAdapter struct {
	client   *cognitoidentityprovider.CognitoIdentityProvider
	PoolID   string
	ClientID string
}

func NewCognitoAdapter(region, poolID string, clientID string) (*CognitoAdapter, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	return &CognitoAdapter{
		client:   cognitoidentityprovider.New(sess),
		PoolID:   poolID,
		ClientID: clientID,
	}, nil
}

func (a *CognitoAdapter) Register(ctx context.Context, userRegisteration entity.UserRegistration) (string, error) {
	attributes := []*cognitoidentityprovider.AttributeType{
		{
			Name:  aws.String("name"),
			Value: aws.String(userRegisteration.Name),
		},
		{
			Name:  aws.String("email"),
			Value: aws.String(userRegisteration.Email),
		},
	}

	input := &cognitoidentityprovider.SignUpInput{
		ClientId:       aws.String(a.ClientID),
		Password:       aws.String(userRegisteration.Password),
		Username:       aws.String(userRegisteration.Email),
		UserAttributes: attributes,
	}

	result, err := a.client.SignUp(input)
	if err != nil {
		return "", fmt.Errorf("failed to sign up user: %w", err)
	}

	return *result.CodeDeliveryDetails.Destination, nil
}

func (a *CognitoAdapter) Login(ctx context.Context, userLogin entity.UserLogin) (*entity.LoginResult, error) {
	params := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(userLogin.Email),
			"PASSWORD": aws.String(userLogin.Password),
		},
		ClientId: aws.String(a.ClientID),
	}

	result, err := a.client.InitiateAuth(params)
	if err != nil {
		return nil, fmt.Errorf("failed to initate auth: %w", err)
	}

	loginReturn := entity.LoginResult{
		AccessToken:  result.AuthenticationResult.AccessToken,
		ExpiresIn:    result.AuthenticationResult.ExpiresIn,
		IdToken:      result.AuthenticationResult.IdToken,
		RefreshToken: result.AuthenticationResult.RefreshToken,
		TokenType:    result.AuthenticationResult.TokenType,
	}

	return &loginReturn, nil
}

func (a *CognitoAdapter) GetUser(ctx context.Context, email string) (*cognitoidentityprovider.AdminGetUserOutput, error) {
	input := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(a.PoolID),
		Username:   aws.String(email),
	}

	result, err := a.client.AdminGetUser(input)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return result, nil
}

func (a *CognitoAdapter) ConfirmRegistration(ctx context.Context, userRegistrationConfirm entity.UserRegistrationConfirm) error {
	input := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(a.ClientID),
		ConfirmationCode: aws.String(userRegistrationConfirm.ConfirmationCode),
		Username:         aws.String(userRegistrationConfirm.Email),
	}

	_, err := a.client.ConfirmSignUp(input)
	if err != nil {
		return fmt.Errorf("failed to confirm sign up: %w", err)
	}

	return nil
}

func (a *CognitoAdapter) ResendConfirmationCode(ctx context.Context, email entity.Email) (string, error) {
	input := &cognitoidentityprovider.ResendConfirmationCodeInput{
		ClientId: aws.String(a.ClientID),
		Username: aws.String(email.Email),
	}

	result, err := a.client.ResendConfirmationCode(input)
	if err != nil {
		return "", fmt.Errorf("failed to resend confirmation code: %w", err)
	}

	return *result.CodeDeliveryDetails.Destination, nil
}

func (a *CognitoAdapter) ForgotPassword(ctx context.Context, email entity.Email) (string, error) {
	input := &cognitoidentityprovider.ForgotPasswordInput{
		ClientId: aws.String(a.ClientID),
		Username: aws.String(email.Email),
	}

	result, err := a.client.ForgotPassword(input)
	if err != nil {
		return "", fmt.Errorf("failed to initiate forgot password: %w", err)
	}

	return *result.CodeDeliveryDetails.Destination, nil
}

func (a *CognitoAdapter) ConfirmForgotPassword(ctx context.Context, userResetPassword entity.UserResetPassword) (*cognitoidentityprovider.ConfirmForgotPasswordOutput, error) {
	input := &cognitoidentityprovider.ConfirmForgotPasswordInput{
		ClientId:         aws.String(a.ClientID),
		Username:         aws.String(userResetPassword.Email),
		ConfirmationCode: aws.String(userResetPassword.ConfirmationCode),
		Password:         aws.String(userResetPassword.NewPassword),
	}

	result, err := a.client.ConfirmForgotPassword(input)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm forgot password: %w", err)
	}

	return result, err
}

func (a *CognitoAdapter) ChangePassword(ctx context.Context, accessToken string, userChangePassword entity.UserChangePassword) error {
	input := &cognitoidentityprovider.ChangePasswordInput{
		AccessToken:      aws.String(accessToken),
		PreviousPassword: aws.String(userChangePassword.PreviousPassword),
		ProposedPassword: aws.String(userChangePassword.ProposedPassword),
	}

	_, err := a.client.ChangePassword(input)
	if err != nil {
		return fmt.Errorf("failed to change password: %w", err)
	}

	return nil
}
