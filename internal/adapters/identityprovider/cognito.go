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
	},nil
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

func (a *CognitoAdapter) Login(ctx context.Context, userLogin entity.UserLogin) (string, error) {
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
		return "", fmt.Errorf("failed to initate auth: %w", err)
	}

	return *result.AuthenticationResult.IdToken, nil
}
