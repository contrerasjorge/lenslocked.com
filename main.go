package main

import (
	"fmt"
	"net/http"

	"lenslocked.com/controllers"
	"lenslocked.com/models"

	"github.com/gorilla/mux"
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
	must(err)
	defer us.Close()
	// us.DestructiveReset()
	us.AutoMigrate()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(us)

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/signup", usersC.NewView).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
