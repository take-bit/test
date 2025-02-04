package authUseCase

import (
	"context"
	"github.com/take-bit/auth-service/internal/domain/auth/dto"
)

type LoginUserUseCase struct {

}

func NewLoginUserUseCase() *LoginUserUseCase {
	return &LoginUserUseCase{}
}

func (l *LoginUserUseCase) Execute(ctx context.Context, req *dto.LoginUserRequest) (*dto.LoginUserResponse, error) {
	panic("todo")
}
