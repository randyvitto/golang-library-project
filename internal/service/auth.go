package service

import (
	"belajar-golang-rest-api/lat/domain"
	"belajar-golang-rest-api/lat/dto"
	"belajar-golang-rest-api/lat/internal/config"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authservice struct {
	conf           *config.Config
	UserRepository domain.UserRepository
}



func NewAuth(cnf *config.Config,
	userRepository domain.UserRepository) domain.AuthService {
	return authservice{
		conf:           cnf,
		UserRepository: userRepository,
	}
}

// Login implements domain.AuthService.
func (a authservice) Login(ctx context.Context, req dto.AuthRequest) (dto.AuthResponse, error) {
	user, err := a.UserRepository.FindByEmail(ctx, req.Email)
	if err != nil{
		return dto.AuthResponse{}, err
	}
	if user.Id == ""{
		return dto.AuthResponse{} , errors.New("Authentication Failed")	
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))  
	if err != nil{
		return dto.AuthResponse{} , errors.New("Authentication Failed")
	}
	claim := jwt.MapClaims{
		"id" : user.Id,
		"exp": time.Now().Add(time.Duration(a.conf.Jwt.Exp) * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	TokenStr, err := token.SignedString([]byte(a.conf.Jwt.Key))
	if err != nil{
		return dto.AuthResponse{} , errors.New("Authentication Failed")
	}
	return dto.AuthResponse{
		Token: TokenStr,
	}, nil

}