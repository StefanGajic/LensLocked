package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	//ErrNotFound is returned when a resource cannot be found
	//in the database
	ErrNotFound = errors.New("models: resources not found")

	//ErrInvalidID is returned when an invalid ID is provided
	//to a method like Delete
	ErrInvalidID = errors.New("models: ID provided was invalid")
)

func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &UserService{
		db: db,
	}, nil
}

type UserService struct {
	db *gorm.DB
}

//ByID will look up the id provided
//1 case- user, nil
//2 case- nil, ErrNotFound
//3 case- nil, otherError
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err

}

//ByEmail looks up a user with the given email address and
//returns that user.

func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

//first will query using the provided gorm.DB and it will
//getthe first item returnetand place it into destination
//Id nothing is found in query it will return ErrNotFound
func first(db *gorm.DB, destination interface{}) error {
	err := db.First(destination).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

//Create the provited user and backfill data like:
//ID, CreatedAt and UpdatedAt fields.
func (us *UserService) Create(user *User) error {
	return us.db.Create(user).Error
}

//Update will update provided user with all of the date
//inthe provided user object
func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

//will delete the user with the provided ID
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}

//Closes the UserService database connection
func (us *UserService) Close() error {
	return us.db.Close()
}

//DestructiveReset drops the user table and rebuilds it
func (us *UserService) DestructiveReset() error {
	if err := us.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return us.AutoMigrate()
}

//AutoMigrate will attempt to automatically
//migrate the users table
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}
