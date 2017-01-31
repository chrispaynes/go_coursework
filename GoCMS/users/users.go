package users

import (
	"errors"
	"fmt"
	"time"

	"github.com/boltdb/bolt"

	"golang.org/x/crypto/bcrypt"
)

// DB represents an in-memory user database.
var DB = newDB()

// ErrUserAlreadyExists represents a error that is thrown when a user attempts to register with a username that is already present in the database.
var ErrUserAlreadyExists = errors.New("User Already Exists!")

// ErrUserNotFound indicates a user was not found within the database.
var ErrUserNotFound = errors.New("User Not Found")

// Store represents an in-memory database to store user information. The contains two internal stores: Users and a Session store.
type Store struct {
	DB       *bolt.DB
	Users    string
	Sessions string
}

// newDB construct a BoltDB data store and initializes buckets.
func newDB() *Store {
	db, err := bolt.Open("users.db", 0600, &bolt.Options{
		Timeout: 1 * time.Second,
	})

	if err != nil {
		panic(err)
	}

	s := &Store{
		DB:       db,
		Users:    "Users",
		Sessions: "Sessions",
	}

	s.createBucket("Users")
	s.createBucket("Sessions")

	return s
}

// createBucket creates a BoltDB bucket to store Key/Value pair collections.
func (s *Store) createBucket(name string) error {
	s.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	return nil
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

	if hashedPassword == nil {
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
