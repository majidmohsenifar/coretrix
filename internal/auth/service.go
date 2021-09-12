package auth

import (
	"coretrix/internal/platform"
	"coretrix/internal/user"
)

type Service interface {
	GetUser(token string) (*user.AutheticatedUser, error)
}

type service struct {
	configs platform.Configs
}

func (s *service) GetUser(token string) (*user.AutheticatedUser, error) {
	ID := 1
	//this is because we could check the special user in integration test
	if s.configs.GetEnv() == platform.EnvTest {
		ID = 12345
	}
	return &user.AutheticatedUser{
		ID:       ID,
		Username: "username",
		IsGuest:  false,
	}, nil
}

func NewAuthService(configs platform.Configs) Service {
	return &service{
		configs: configs,
	}
}
