package users

import (
	"errors"
	"time"

	"github.com/boltdb/bolt"

	"golang.org/x/crypto/bcrypt"
)

// Store represents an in-memory database to store user information. Uses Synchronization Mutex locks to prevent multiple goroutines from modifying the map at the same time.
type Store struct {
	DB       *bolt.DB
	Users    string
	Sessions string
}

// ErrUserAlreadyExists represents a error that is thrown when a user attempts to register with a username that is already present in the database.
var ErrUserAlreadyExists = errors.New("User Already Exists!")
var ErrUserNotFound = errors.New("User Not Found")

// DB represents an in-memory user database.
var DB = newDB()

func newDB() *Store {
	db, err := bolt.Open("users.db", 0600, &bolt.Options{
		Timeout: 1 * time.Second,
	})

	if err != nil {
		panic(err)
	}

	return &Store{
		DB:       db,
		Users:    "Users",
		Sessions: "Sessions",
	}
}

func (s *Store) createBucket(name string) error {
	var err error
	s.DB.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("Users"))
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// NewUser registers and stores a new user with a username and hashed password in the database.
func NewUser(username, password string) error {
	err := exists(username)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return DB.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DB.Users))
		return b.Put([]byte(username), hashedPassword)
	})
}

// AuthenticateUser verifies the username and password matches the username and hashed password stored in the database.
func AuthenticateUser(username, password string) error {
	var hashedPassword []byte
	DB.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DB.Users))
		hashedPassword = b.Get([]byte(username))
		return nil
	})
	if hashedPassword != nil {
		return ErrUserNotFound
	}

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// OverrideOldPassword overrides a stored password with a new password.
func OverrideOldPassword(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return DB.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DB.Users))
		return b.Put([]byte(username), hashedPassword)
	})
}

// exists checks if a username already exists within the database.
func exists(username string) error {
	var result []byte
	DB.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DB.Users))
		result = b.Get([]byte(username))
		return nil
	})
	if result != nil {
		return ErrUserAlreadyExists
	}
	return nil
}
