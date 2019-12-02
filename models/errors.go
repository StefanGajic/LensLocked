package models

import "strings"

const (
	//ErrNotFound is returned when a resource cannot be found
	//in the database
	ErrNotFound modelError = "models: resources not found"

	ErrInvalidPassword modelError = "models: incorrect password provided"

	ErrEmailRequired modelError = "models: Email address is required"

	ErrEmailInvalid modelError = "models: Email address is not valid"

	ErrEmailTaken modelError = "models: email address is already taken"

	ErrPasswordTooShort modelError = "models: password must be at least 8 characters long"

	ErrTitleRequired modelError = "models: title is required"

	ErrPasswordRequired modelError = "models: password is required"

	//ErrInvalidID is returned when an invalid ID is provided
	//to a method like Delete
	ErrInvalidID privateError = "models: ID provided was invalid"

	ErrRememberTooShort privateError = "models: remember token must be at least 32 bytes"

	ErrRememberRequired privateError = "models: remember token is required"

	ErrUserIDRequired privateError = "models: user ID is required"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}

type privateError string

func (e privateError) Error() string {
	return string(e)
}
