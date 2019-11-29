package models

import "strings"

const (
	//ErrNotFound is returned when a resource cannot be found
	//in the database
	ErrNotFound modelError = "models: resources not found"

	//ErrInvalidID is returned when an invalid ID is provided
	//to a method like Delete
	ErrInvalidID modelError = "models: ID provided was invalid"

	ErrInvalidPassword modelError = "models: incorrect password provided"

	ErrEmailRequired modelError = "models: Email address is required"

	ErrEmailInvalid modelError = "models: Email address is not valid"

	ErrEmailTaken modelError = "models: email address is already taken"

	ErrRememberTooShort modelError = "models: remember token must be at lear 32 bytes"

	ErrPasswordRequired modelError = "models: password is required"

	ErrPasswordTooShort modelError = "models: password must be at least 8 characters long"

	ErrRememberRequired modelError = "models: remember token is required"
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
