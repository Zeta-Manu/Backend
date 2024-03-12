package entity

type UserRegistration struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegistrationConfirm struct {
	Email            string `json:"email"`
	ConfirmationCode string `json:"confirmation_code"`
}

type UserResetPassword struct {
	Email            string `json:"email"`
	ConfirmationCode string `json:"confirmation_code"`
	NewPassword      string `json:"new_password"`
}

type UserChangePassword struct {
	PreviousPassword string `json:"previous_password"`
	ProposedPassword string `json:"proposed_password"`
}

type Email struct {
	Email string `json:"email"`
}

type LoginResult struct {
	AccessToken  *string `json:"access_token"`
	ExpiresIn    *int64  `json:"expires_in"`
	IdToken      *string `json:"id_token"`
	RefreshToken *string `json:"refresh_token"`
	TokenType    *string `json:"token_type"`
}
