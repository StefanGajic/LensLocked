package controllers

import (
	"net/http"
	"time"

	"github.com/lenslocked/context"
	"github.com/lenslocked/models"
	"github.com/lenslocked/rand"

	"github.com/lenslocked/views"
)

//NewsUsers is used to create a new Users controller
//this func will panic if templates are not
//parsed correctly, and should only be used
//during initial setup
func NewUsers(us models.UserService) *Users { //, emailer *email.Client
	return &Users{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
		//emailer:   emailer,
	}
}

type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        models.UserService
}

// new is used to render the form where a user can
//create a new user acc

// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	parseURLParams(r, &form)
	u.NewView.Render(w, r, form)
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
	var vd views.Data
	var form SignupForm
	vd.Yield = &form
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}
	//   u.emailer.Welcome(user.Name, user.Email)
	err := u.signIn(w, &user)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return

	}
	alert := views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Welcome to Lens Gopher!!!",
	}
	views.RedirectAlert(w, r, "/galleries", http.StatusFound, alert)

}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

//used to verify provided email and passand then log the
//user in if they are correct
//POST / login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}
	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			vd.AlertError("Invalid email address")
		default:
			vd.SetAlert(err)
		}
		u.LoginView.Render(w, r, vd)
		return
	}

	err = u.signIn(w, user)
	if err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

// POST /logout
func (u *Users) Logout(w http.ResponseWriter, r *http.Request) {
	// First expire the user's cookie
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	// Then we update the user with a new remember token
	user := context.User(r.Context())
	// We are ignoring errors for now because they are
	// unlikely, and even if they do occur we can't recover
	// now that the user doesn't have a valid cookie
	token, _ := rand.RememberToken()
	user.Remember = token
	u.us.Update(user)
	// Finally send the user to the home page
	http.Redirect(w, r, "/", http.StatusFound)
}

// used to sign in via cookies
func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}

	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	return nil
}
