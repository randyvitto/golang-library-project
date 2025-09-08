package domain

import (
	"belajar-golang-rest-api/lat/dto"
	"context"
)

type AuthService interface {
	Login(ctx context.Context, req dto.AuthRequest) (dto.AuthResponse, error)
}