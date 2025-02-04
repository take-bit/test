package dto

type LoginUserRequest struct {
	Login    string
	Password string
	Type     string
}

type LoginUserResponse struct {
	AccessToken  string
	RefreshToken string
}
