package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/csrf"
	"golang.org/x/oauth2"
	llctx "lenslocked.com/context"
	"lenslocked.com/models"
)

// NewOAuths is used to create a new OAuths controller,
// This function will panic if the templates are not
// parsed correctly, and should only be used during
// initial setup
func NewOAuths(os models.OAuthService, dbxConfig *oauth2.Config) *OAuths {
	return &OAuths{
		os:       os,
		dbxOAuth: dbxConfig,
	}
}

type OAuths struct {
	os       models.OAuthService
	dbxOAuth *oauth2.Config
}

func (o *OAuths) DropboxConnect(w http.ResponseWriter, r *http.Request) {
	state := csrf.Token(r)
	cookie := http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	url := o.dbxOAuth.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

func (o *OAuths) DropboxCallback(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	state := r.FormValue("state")
	cookie, err := r.Cookie("oauth_state")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if cookie == nil || cookie.Value != state {
		http.Error(w, "Invalid state provided", http.StatusBadRequest)
		return
	}
	cookie.Value = ""
	cookie.Expires = time.Now()
	http.SetCookie(w, cookie)

	code := r.FormValue("code")
	token, err := o.dbxOAuth.Exchange(context.TODO(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := llctx.User(r.Context())
	existing, err := o.os.Find(user.ID, models.OAuthDropbox)
	if err == models.ErrNotFound {
		// noop
	} else if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		o.os.Delete(existing.ID)
	}
	userOAuth := models.OAuth{
		UserID:  user.ID,
		Token:   *token,
		Service: models.OAuthDropbox,
	}
	err = o.os.Create(&userOAuth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "%+v", token)
}
