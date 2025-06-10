package database

import (
	"errors"
	"sync"

	"github.com/TRu-S3/backend/internal/domain/entity"
)

type UserRepositoryImpl struct {
	users  map[int]*entity.User
	nextID int
	mu     sync.RWMutex
}

func NewUserRepositoryImpl() *UserRepositoryImpl {
	return &UserRepositoryImpl{
		users:  make(map[int]*entity.User),
		nextID: 1,
	}
}

func (r *UserRepositoryImpl) Create(user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user.ID = r.nextID
	r.users[r.nextID] = user
	r.nextID++

	return nil
}

func (r *UserRepositoryImpl) GetByID(id int) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (r *UserRepositoryImpl) GetAll() ([]*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*entity.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepositoryImpl) Update(user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return errors.New("user not found")
	}

	r.users[user.ID] = user
	return nil
}

func (r *UserRepositoryImpl) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return errors.New("user not found")
	}

	delete(r.users, id)
	return nil
}
