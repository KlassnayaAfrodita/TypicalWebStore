package storage

import (
	"sync"
)

// var NotFound = errors.New("not found")

type User struct {
	ID       int             `json:"user_id"`
	Name     string          `json:"user_name"`
	Password string          `json:"user_password"`
	Cart     *ProductStorage `json:"user_cart"`
}

type UserStorage struct {
	users  map[int]User
	mu     *sync.RWMutex
	nextID int
}

func NewUserStorage() *UserStorage {
	return &UserStorage{
		users: map[int]User{},
		mu:    &sync.RWMutex{},
	}
}

func (us *UserStorage) GetUserByName(name string) (User, error) {
	for _, user := range us.users {
		if user.Name == name {
			return user, nil
		}
	}
	return User{}, NotFound
}

func (us *UserStorage) AddUser(in User) (User, error) {
	us.mu.Lock()
	defer us.mu.Unlock()

	us.nextID++
	in.ID = us.nextID
	us.users[in.ID] = in

	return in, nil
}

func (us *UserStorage) GetUser(id int) (User, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	user, ok := us.users[id]
	if !ok {
		return User{}, NotFound
	}
	return user, nil
}

func (us *UserStorage) GetUsers() ([]User, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	result := make([]User, 0, len(us.users))
	for _, user := range us.users {
		result = append(result, user)
	}

	return result, nil
}

func (us *UserStorage) ChangeUser(newUser User) (User, error) {
	_, ok := us.users[newUser.ID]

	if !ok {
		return User{}, NotFound
	}

	us.users[newUser.ID] = newUser

	return newUser, nil
}
