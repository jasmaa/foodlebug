// Handlers for account login and creation
package foodlebug

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jasmaa/foodlebug/internal/auth"
	"github.com/jasmaa/foodlebug/internal/store"
)

// Account creation page handler
func handleCreateAccount(store *store.Store) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		messages := make([]string, 10)

		switch r.Method {
		case "GET":
			displayPage(w, "assets/templates/createAccount.html", false, messages)

		case "POST":
			username := r.FormValue("username")
			email := r.FormValue("email")
			password := r.FormValue("password")
			confirm := r.FormValue("confirm-password")

			if username == "" || password == "" || email == "" {
				messages := append(messages, "Fields cannot be empty.")
				displayPage(w, "assets/templates/createAccount.html", false, messages)
				return
			} else if password != confirm {
				messages := append(messages, "Passwords do not match.")
				displayPage(w, "assets/templates/createAccount.html", false, messages)
				return
			}

			err := auth.CreateNewUser(store, username, email, password)

			if err != nil {
				messages := append(messages, "Account could not be created.")
				displayPage(w, "assets/templates/createAccount.html", false, messages)
				return
			}

			// go to success page
			displayPage(w, "assets/templates/createAccountSuccess.html", false, "")
		}
	})
}

// Login page handler
func handleLogin(store *store.Store) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		messages := make([]string, 10)

		switch r.Method {

		case "GET":
			displayPage(w, "assets/templates/login.html", false, messages)

		case "POST":
			username := r.FormValue("username")
			password := r.FormValue("password")

			user, err := store.GetUser("username", username)
			if err != nil {
				messages := append(messages, "User does not exist.")
				displayPage(w, "assets/templates/login.html", false, messages)
				return
			}

			if !auth.VerifyPassword(user.Password, password) {
				messages := append(messages, "Incorrect password.")
				displayPage(w, "assets/templates/login.html", false, messages)
				return
			}

			// set session cookie
			s, err := auth.NewSession(store, user, time.Time{}, "ip address", "user agent")
			if err != nil {
				messages := append(messages, "Could not login.")
				displayPage(w, "assets/templates/login.html", false, messages)
				return
			}

			cookie := &http.Cookie{
				Name:     "foodlebug",
				Value:    fmt.Sprintf("%s_%s", s.UserKey, s.SessionId),
				HttpOnly: true,
				Path:     "/",
				Expires:  s.Expires,
			}

			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
		}
	})
}

// Logout handler
func handleLogout(store *store.Store) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.SessionToUser(store, r)

		// Redirect to login on session failure
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		err = auth.ExpireSession(store, user.Username)

		// Redirect to profile on error
		if err != nil {
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
			return
		}

		displayPage(w, "assets/templates/logoutSuccess.html", false, "")
	})
}
