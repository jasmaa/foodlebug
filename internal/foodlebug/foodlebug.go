// Main app
package foodlebug

import (
	"fmt"
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

func (f *Foodlebug) Run() {
	// Runs site
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
