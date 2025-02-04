package authUseCase

import (
	"context"
	"github.com/take-bit/auth-service/internal/domain/auth/dto"
)

type UpdateAccessTokenUseCase struct {

}

func NewUpdateAccessTokenUseCase() *UpdateAccessTokenUseCase {
	return &UpdateAccessTokenUseCase{}
}

func (u *UpdateAccessTokenUseCase) Execute(ctx context.Context, req *dto.UpdateAccessTokenRequest) (*dto.UpdateAccessTokenResponse, error) {
	panic("todo")
}
