package service

import (
	"net/http"
	"tokopedia.se.training/Project1/usermanagesys/api/configuration"
	"tokopedia.se.training/Project1/usermanagesys/api/repository"
)

type UserService struct {
	configuration *configuration.Configuration
	userRepository *repository.TokopediaUserRepository
}

func NewUserService(config *configuration.Configuration,
	userRepo *repository.TokopediaUserRepository) *UserService {
	return &UserService{
		configuration: config,
	}
}

func (s *UserService) GetUsers(w http.ResponseWriter, r *http.Request){

}

func (s *UserService) GetUsersPaging(w http.ResponseWriter, r *http.Request){

}