package users

import (
	"errors"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

// Store represents an in-memory database to store user information. Uses Synchronization Mutex locks to prevent multiple goroutines from modifying the map at the same time.
type Store struct {
	rwm *sync.RWMutex
	m   map[string]string
}

// ErrUserAlreadyExists represents a error that is thrown when a user attempts to register with a username that is already present in the database.
var ErrUserAlreadyExists = errors.New("User Already Exists!")

// DB represents an in-memory user database.
var DB = newDB()

func newDB() *Store {
	return &Store{
		rwm: &sync.RWMutex{},
		m:   make(map[string]string),
	}
}

// NewUser registers and stores a new user with a username and hashed password in the database.
func NewUser(username, password string) error {
	err := exists(username)
	if err != nil {
		return err
	}
	DB.rwm.Lock()
	defer DB.rwm.Unlock()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	DB.m[username] = string(hashedPassword)
	return nil

}

// AuthenticateUser verifies the username and password matches the username and hashed password stored in the database.
func AuthenticateUser(username, password string) error {
	DB.rwm.RLock()
	defer DB.rwm.RUnlock()

	hashedPassword := DB.m[username]
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err
}

// OverrideOldPassword overrides a stored password with a new password.
func OverrideOldPassword(username, password string) error {
	DB.rwm.Lock()
	defer DB.rwm.Unlock()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	DB.m[username] = string(hashedPassword)
	return nil

}

// exists checks if a username already exists within the database.
func exists(username string) error {
	DB.rwm.RLock()
	defer DB.rwm.RUnlock()

	if DB.m[username] != "" {
		return ErrUserAlreadyExists
	}

	return nil
}
