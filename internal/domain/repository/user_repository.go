package repository

import "github.com/TRu-S3/backend/internal/domain/entity"

type UserRepository interface {
	Create(user *entity.User) error
	GetByID(id int) (*entity.User, error)
	GetAll() ([]*entity.User, error)
	Update(user *entity.User) error
	Delete(id int) error
}
