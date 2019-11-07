package controllers

import (
	"fmt"
	"net/http"

	"github.com/lenslocked/models"

	"github.com/lenslocked/views"
)

//NewsUsers is used to create a new Users controller
//this func will panic if templates are not
//parsed correctly, and should only be used
//during initial setup
func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "users/new"),
		us:      us,
	}
}

type Users struct {
	NewView *views.View
	us      *models.UserService
}

// new is used to render the form where a user can
//create a new user acc

// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// create is used to process the signup form when a user
//submit it. This is used to create a new user acc

// POST / signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {

	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user := models.User{
		Name:  form.Name,
		Email: form.Email,
	}
	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, form)

}
