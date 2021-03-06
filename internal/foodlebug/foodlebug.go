// Main app
package foodlebug

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jasmaa/foodlebug/internal/store"
)

const (
	SERVER_PORT = 8000
)

type Foodlebug struct {
	store *store.Store
}

// Runs site
func (f *Foodlebug) Run() {

	fmt.Println("Start...")

	// connect to db
	f.store = &store.Store{}
	f.store.Connect(host, port, user, password, dbname)

	// routing
	r := mux.NewRouter()
	r.Handle("/login", handleLogin(f.store))
	r.Handle("/logout", handleLogout(f.store))
	r.Handle("/createAccount", handleCreateAccount(f.store))

	r.Handle("/postEntry", handlePostEntry(f.store))
	r.Handle("/profile", handleProfile(f.store))
	r.Handle("/about", handleAbout(f.store))
	r.Handle("/nearby", handleNearby(f.store))
	r.Handle("/browse", handleBrowse(f.store))
	r.Handle("/page/{postId}", handlePage(f.store))
	r.Handle("/", handleHome(f.store))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("assets/static"))))
	http.Handle("/", r)

	// start server
	fmt.Printf("Starting server at %d...\n", SERVER_PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%d", SERVER_PORT), nil)
	if err != nil {
		panic(err)
	}
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
		"assets/templates/includes/components.html",
		contentPath,
		navbarPath,
	)
	t.ExecuteTemplate(w, "main", data)
}
