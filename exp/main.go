package main

import (
	"fmt"

	_ "github.com/lib/pq"
	"lenslocked.com/models"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "contrerasjorge"
	dbname = "lenslocked_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()

	us.DestructiveReset()

	user := models.User{
		Name:  "Jimmy Jim",
		Email: "abc@abc.io",
	}
	if err := us.Create(&user); err != nil {
		panic(err)
	}

	// user.Email = "abc@cde.io"
	// if err := us.Update(&user); err != nil {
	// 	panic(err)
	// }

	// UserByEmail, err := us.ByEmail("abc@cde.io")

	if err := us.Delete(user.ID); err != nil {
		panic(err)
	}
	userByID, err := us.ByID(user.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println(userByID)
}
