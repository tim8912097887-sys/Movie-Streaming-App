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

type UserServiceConfig struct {
	PasswordService PasswordService
	JWTService      JWTService
	EnvConfigs      configs.Configs
}

type service struct {
	passwordService PasswordService
	jwtService      JWTService
	envConfigs      configs.Configs
}

func NewUserService(userServiceConfig UserServiceConfig) *service {
	return &service{
		passwordService: userServiceConfig.PasswordService,
		jwtService:      userServiceConfig.JWTService,
		envConfigs:      userServiceConfig.EnvConfigs,
	}
}

var inmemoryUser = []User{}

func (s *service) createUser(ctx context.Context, userPayload CreateUserSchema) (UserDTO, error) {

	for _, user := range inmemoryUser {
		if user.Email == userPayload.Email {
			return UserDTO{}, ErrUserAlreadyExists
		}
	}
	 hashedPassword, err := s.passwordService.HashPassword(userPayload.Password)
	 if err != nil {
		return UserDTO{}, err
	 }
	 user := User{
		ID:        userPayload.Name+userPayload.Email,
		Name:      userPayload.Name,
		Email:     userPayload.Email,
		Password:  hashedPassword,
		Favorite_Genres: userPayload.Favorite_Genres,
	 }
	 inmemoryUser = append(inmemoryUser, user)
	 return s.DataToDTO(user), nil
}

func (s *service) loginUser(ctx context.Context, userPayload LoginUserSchema) (TokenResponse, error) {
	
    var existingUser User

	for _, user := range inmemoryUser {
		if user.Email == userPayload.Email {
			existingUser = user
		}
	}

	if existingUser.ID == "" {
		return TokenResponse{}, ErrInvalidCredentials
	}

	if !s.passwordService.CheckPasswordHash(userPayload.Password, existingUser.Password) {
		return TokenResponse{}, ErrInvalidCredentials
	}

	refreshTokenPayload := auth.JWTGeneratePayload{
		Subject: existingUser.ID,
		TokenVersion: existingUser.TokenVersion,
		Duration: RefreshTokenExpiredTime,
		Secret: s.envConfigs.RefreshTokenSecret,
	}

	accessTokenPayload := auth.JWTGeneratePayload{
		Subject: existingUser.ID,
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
	for i, user := range inmemoryUser {
		if user.ID == userId {
			if user.TokenVersion != tokenVersion {
				return TokenResponse{}, ErrTokenVersionMismatch
			}
			inmemoryUser[i].TokenVersion = user.TokenVersion + 1
			refreshTokenPayload := auth.JWTGeneratePayload{
				Subject: inmemoryUser[i].ID,
				TokenVersion: inmemoryUser[i].TokenVersion,
				Duration: RefreshTokenExpiredTime,
				Secret: s.envConfigs.RefreshTokenSecret,
			}

			accessTokenPayload := auth.JWTGeneratePayload{
				Subject: inmemoryUser[i].ID,
				TokenVersion: inmemoryUser[i].TokenVersion,
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
			return TokenResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
		}
	}
	return TokenResponse{}, ErrUserNotFound
}

func (s *service) logoutUser(ctx context.Context, userId string,tokenVersion int) error {
	
	for i, user := range inmemoryUser {
		if user.ID == userId {
			if user.TokenVersion != tokenVersion {
				return ErrTokenVersionMismatch
			}
			inmemoryUser[i].TokenVersion = user.TokenVersion + 1
			return nil
		}
	}
	return ErrUserNotFound
}

// For debug
func (s *service) getAllUsers(ctx context.Context) ([]UserDTO, error) {
	response := make([]UserDTO, 0, len(inmemoryUser))
	for _, user := range inmemoryUser {
		response = append(response, s.DataToDTO(user))
	}
	return response, nil
}

func (s *service) DataToDTO(user User) UserDTO {
	return UserDTO{
		Name:      user.Name,
		Email:     user.Email,
	}
}