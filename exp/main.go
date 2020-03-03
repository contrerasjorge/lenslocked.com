package main

import (
	"fmt"

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
	us, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.AutoMigrate()

	User := models.User{
		Name:     "jei sea",
		Email:    "jei@jei.com",
		Password: "jei",
		Remember: "abc123",
	}
	err = us.Create(&User)
	if err != nil {
		panic(err)

	}

	user2, err := us.ByRemember("abc123")
	if err != nil {
		panic(err)
	}
	fmt.Println(user2)
}
