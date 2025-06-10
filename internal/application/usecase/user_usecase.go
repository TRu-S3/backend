package usecase

import (
	"time"

	"github.com/TRu-S3/backend/internal/application/dto"
	"github.com/TRu-S3/backend/internal/domain/entity"
	"github.com/TRu-S3/backend/internal/domain/repository"
)

type UserUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

func (uc *UserUseCase) CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	user := &entity.User{
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (uc *UserUseCase) GetUser(id int) (*dto.UserResponse, error) {
	user, err := uc.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (uc *UserUseCase) GetAllUsers() ([]*dto.UserResponse, error) {
	users, err := uc.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.UserResponse, len(users))
	for i, user := range users {
		responses[i] = &dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}
	}

	return responses, nil
}

func (uc *UserUseCase) UpdateUser(id int, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := uc.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	user.UpdatedAt = time.Now()

	if err := uc.userRepo.Update(user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (uc *UserUseCase) DeleteUser(id int) error {
	return uc.userRepo.Delete(id)
}
