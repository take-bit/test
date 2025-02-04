package dto

type UpdateAccessTokenRequest struct {
	RefreshToken string
}
type UpdateAccessTokenResponse struct {
	AccessToken string
}
