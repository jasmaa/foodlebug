package foodlebug

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jasmaa/foodlebug/internal/models"

	"github.com/jasmaa/foodlebug/internal/auth"
	"github.com/jasmaa/foodlebug/internal/store"
)

// Handle home page
func handleHome(store *store.Store) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		_, err := auth.SessionToUser(store, r)
		posts, _ := store.GetPosts()

		if err != nil {
			displayPage(w, "assets/templates/home.html", false, posts)
			return
		}

		displayPage(w, "assets/templates/home.html", true, posts)
	})
}

// Handle about page
func handleAbout(store *store.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := auth.SessionToUser(store, r)

		if err != nil {
			displayPage(w, "assets/templates/about.html", false, "")
			return
		}

		displayPage(w, "assets/templates/about.html", true, "")
	})
}

// Handle post page
func handlePage(store *store.Store) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		post, err := store.GetPost("id", vars["postId"])
		if err != nil || post == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// check if logged in
		_, err = auth.SessionToUser(store, r)
		if err != nil {
			displayPage(w, "assets/templates/postPage.html", false, post)
			return
		}

		displayPage(w, "assets/templates/postPage.html", true, post)
	})
}

// Post entry page handler
func handlePostEntry(store *store.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user, err := auth.SessionToUser(store, r)

		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		messages := make([]string, 10)

		switch r.Method {

		case "GET":
			displayPage(w, "assets/templates/postEntry.html", true, messages)

		case "POST":
			entryTitle := r.FormValue("entry-title")
			content := r.FormValue("entry-content")
			entryPhoto := r.FormValue("entry-photo")

			lat, err := strconv.ParseFloat(r.FormValue("lat"), 64)
			if err != nil {
				messages := append(messages, "Could not submit post.")
				displayPage(w, "assets/templates/postEntry.html", false, messages)
				return
			}

			lon, err := strconv.ParseFloat(r.FormValue("lon"), 64)
			if err != nil {
				messages := append(messages, "Could not submit post.")
				displayPage(w, "assets/templates/postEntry.html", false, messages)
				return
			}

			// Check if entries filled
			if entryTitle == "" || content == "" || entryPhoto == "" {
				messages := append(messages, "Could not submit post.")
				displayPage(w, "assets/templates/postEntry.html", false, messages)
				return
			}

			// Get b64 data
			data := strings.Split(entryPhoto, ",")
			if len(data) != 2 {
				messages := append(messages, "Could not submit post.")
				displayPage(w, "assets/templates/postEntry.html", false, messages)
				return
			}

			// Limit size
			if float64(len(data[1]))/1.37 > 5000000 {
				messages := append(messages, "Image exceeded 500KB.")
				displayPage(w, "assets/templates/postEntry.html", false, messages)
				return
			}

			// Insert post into db
			post := &models.Post{
				Id:          store.GenerateId("posts"),
				PosterId:    user.Id,
				Photo:       data[1],
				Title:       entryTitle,
				Content:     content,
				TimePosted:  time.Now(),
				LocationLat: lat,
				LocationLon: lon,
				Visible:     true,
			}

			err = store.InsertPost(post)
			if err != nil {
				messages := append(messages, "Could not submit post.")
				displayPage(w, "assets/templates/postEntry.html", false, messages)
				return
			}

			// Redirect to home
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

	})
}

// Handles user profile
func handleProfile(store *store.Store) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user, err := auth.SessionToUser(store, r)

		// Redirect to login on session failure
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		displayPage(w, "assets/templates/profile.html", true, user)
	})
}

// Account creation page handler
func handleCreateAccount(store *store.Store) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		messages := make([]string, 10)

		switch r.Method {
		case "GET":
			displayPage(w, "assets/templates/createAccount.html", false, messages)

		case "POST":
			username := r.FormValue("username")
			password := r.FormValue("password")
			confirm := r.FormValue("confirm-password")

			if username == "" || password == "" {
				messages := append(messages, "Fields cannot be empty.")
				displayPage(w, "assets/templates/createAccount.html", false, messages)
				return
			} else if password != confirm {
				messages := append(messages, "Passwords do not match.")
				displayPage(w, "assets/templates/createAccount.html", false, messages)
				return
			}

			err := auth.CreateNewUser(store, username, password)

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

// Display page
func displayPage(w http.ResponseWriter, contentPath string, loggedIn bool, data interface{}) {

	// Select proper navbar
	navbarPath := "assets/templates/includes/navbarLoggedOut.html"
	if loggedIn {
		navbarPath = "assets/templates/includes/navbarLoggedIn.html"
	}

	w.WriteHeader(http.StatusOK)
	var t *template.Template
	t, _ = template.ParseFiles(
		"assets/templates/main.html",
		"assets/templates/includes/footer.html",
		contentPath,
		navbarPath,
	)
	t.ExecuteTemplate(w, "main", data)
}
