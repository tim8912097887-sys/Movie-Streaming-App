package users

import (
	"context"

	"github.com/tim8912097887-sys/server/internal/auth"
	"github.com/tim8912097887-sys/server/internal/configs"
)

type PasswordService interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type JWTService interface {
	GenerateToken(payload auth.JWTGeneratePayload) (string, error)
    ValidateToken(payload auth.JWTValidatePayload) (auth.CustomClaims, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, user CreateUserSchema) (User, error)
	FindUserByEmail(ctx context.Context, email string) (User, error)
	FindUserById(ctx context.Context, id string) (User, error)
	UpdateUser(ctx context.Context, user UpdateUserSchema) error
}

type UserServiceConfig struct {
	Repository      UserRepository
	PasswordService PasswordService
	JWTService      JWTService
	EnvConfigs      configs.Configs
}

type service struct {
	repository      UserRepository
	passwordService PasswordService
	jwtService      JWTService
	envConfigs      configs.Configs
}

func NewUserService(userServiceConfig UserServiceConfig) *service {
	return &service{
		passwordService: userServiceConfig.PasswordService,
		jwtService:      userServiceConfig.JWTService,
		envConfigs:      userServiceConfig.EnvConfigs,
		repository:      userServiceConfig.Repository,
	}
}

func (s *service) createUser(ctx context.Context, userPayload CreateUserSchema) (UserDTO, error) {

	existUser, err := s.repository.FindUserByEmail(ctx, userPayload.Email)
	if err != nil && err != ErrUserNotFound {
		return UserDTO{}, err
	}

	if existUser.Email == userPayload.Email {
		return UserDTO{}, ErrUserAlreadyExists
	}

	 hashedPassword, err := s.passwordService.HashPassword(userPayload.Password)
	 if err != nil {
		return UserDTO{}, err
	 }

	 userPayload.Password = hashedPassword

	 var createdUser User

	 createdUser, err = s.repository.CreateUser(ctx, userPayload)
	 if err != nil {
		return UserDTO{}, err
	 }
	 return s.DataToDTO(createdUser), nil
}

func (s *service) loginUser(ctx context.Context, userPayload LoginUserSchema) (TokenResponse, error) {
	
    var existingUser User

	// check if user exists
	existingUser, err := s.repository.FindUserByEmail(ctx, userPayload.Email)
	if err != nil {
		if err == ErrUserNotFound {
			return TokenResponse{}, ErrInvalidCredentials
		}
		return TokenResponse{}, err
	}

	if !s.passwordService.CheckPasswordHash(userPayload.Password, existingUser.Password) {
		return TokenResponse{}, ErrInvalidCredentials
	}

	// generate tokens
	refreshTokenPayload := auth.JWTGeneratePayload{
		Subject: existingUser.ID.Hex(),
		TokenVersion: existingUser.TokenVersion,
		Duration: RefreshTokenExpiredTime,
		Secret: s.envConfigs.RefreshTokenSecret,
	}

	accessTokenPayload := auth.JWTGeneratePayload{
		Subject: existingUser.ID.Hex(),
		TokenVersion: existingUser.TokenVersion,
		Duration: AccessTokenExpiredTime,
		Secret: s.envConfigs.AccessTokenSecret,
	}

	refreshToken, err := s.jwtService.GenerateToken(refreshTokenPayload)
	if err != nil {
		return TokenResponse{}, err
	}

	accessToken, err := s.jwtService.GenerateToken(accessTokenPayload)
	if err != nil {
		return TokenResponse{}, err
	}

	return TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) refreshToken(ctx context.Context, userId string,tokenVersion int) (TokenResponse, error) {
	
	// check if user exists and token version matches
	existingUser, err := s.repository.FindUserById(ctx, userId)
	if err != nil {
		if err == ErrUserNotFound {
			return TokenResponse{}, ErrUserNotFound
		}
		return TokenResponse{}, err
	}

	if existingUser.TokenVersion != tokenVersion {
		return TokenResponse{}, ErrTokenVersionMismatch
	}

	// update token version
	newTokenVersion := existingUser.TokenVersion + 1
	updateUserSchema := UpdateUserSchema{
		Id:              existingUser.ID.Hex(),
		TokenVersion:    newTokenVersion,
	}
	s.repository.UpdateUser(ctx,updateUserSchema)

	// generate new tokens
	refreshTokenPayload := auth.JWTGeneratePayload{
		Subject: existingUser.ID.Hex(),
		TokenVersion: newTokenVersion,
		Duration: RefreshTokenExpiredTime,
		Secret: s.envConfigs.RefreshTokenSecret,
	}

	accessTokenPayload := auth.JWTGeneratePayload{
		Subject: existingUser.ID.Hex(),
		TokenVersion: newTokenVersion,
		Duration: AccessTokenExpiredTime,
		Secret: s.envConfigs.AccessTokenSecret,
	}

	refreshToken, err := s.jwtService.GenerateToken(refreshTokenPayload)
	if err != nil {
		return TokenResponse{}, err
	}

	accessToken, err := s.jwtService.GenerateToken(accessTokenPayload)
	if err != nil {
		return TokenResponse{}, err
	}

	return TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) logoutUser(ctx context.Context, userId string,tokenVersion int) error {
	
	// check if user exists and token version matches
	existingUser, err := s.repository.FindUserById(ctx, userId)
	if err != nil {
		if err == ErrUserNotFound {
			return ErrUserNotFound
		}
		return err
	}

	if existingUser.TokenVersion != tokenVersion {
		return ErrTokenVersionMismatch
	}

	// update token version
	updateUserSchema := UpdateUserSchema{
		Id:              existingUser.ID.Hex(),
		TokenVersion:    existingUser.TokenVersion + 1,
	}
	s.repository.UpdateUser(ctx,updateUserSchema)

	return nil
}

func (s *service) DataToDTO(user User) UserDTO {
	return UserDTO{
		Name:      user.Name,
		Email:     user.Email,
	}
}