package dto

type RegisterUserResponse struct {
	UUID         string
	AccessToken  string
	RefreshToken string
}

type RegisterUserRequest struct {
	Username string
	Email    string
	Password string
}
