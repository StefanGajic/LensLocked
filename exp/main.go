package main

import (
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lenslocked/models"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "lenslocked_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.DestructiveReset()

	user := models.User{
		Name:  "Rakesh Kumar",
		Email: "rakesh@kumar.in",
	}
	if err := us.Create(&user); err != nil {
		panic(err)
	}
	if err := us.Delete(user.ID); err != nil {
		panic(err)
	}
	// user.Email = "rakesh@india.com"
	// if err := us.Update(&user); err != nil {
	// 	panic(err)
	// }
	// userByEmail, err := us.ByEmail("rakesh@india.com")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(userByEmail)
	userByID, err := us.ByID(user.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println(userByID)
}
