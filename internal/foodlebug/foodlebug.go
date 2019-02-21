// Main app

package foodlebug

import (
  "fmt"
  "net/http"
  "html/template"
  "github.com/gorilla/mux"

  "github.com/jasmaa/foodlebug/internal/store"
)

const (
  SERVER_PORT = 8000
)

type Foodlebug struct {
  store *store.Store
}

func handleDisplay(store *store.Store) http.Handler {
  // Dummy handler

  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    users := store.GetUsers()

    w.WriteHeader(http.StatusOK)
    var t *template.Template
    t, _ = template.ParseFiles(
      "assets/templates/main.html",
      "assets/templates/testDisplay.html",
      "assets/templates/footer.html",
    )
    t.ExecuteTemplate(w, "main", users)
  })
}

func (f *Foodlebug) Run(){
  // Runs site

  fmt.Println("Start...")

  // connect to db
  f.store = &store.Store{}
  f.store.Connect(host, port, user, password, dbname)

  /*
  store.AddUser(&User{
    2,
    "John",
    "super secret password",
    6.9,
  })
  */

  // routing
  r := mux.NewRouter()
  r.Handle("/", handleDisplay(f.store))
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("assets/static"))))
  http.Handle("/", r)

  // start server
  fmt.Printf("Starting server at %d...\n", SERVER_PORT)
  err := http.ListenAndServe(fmt.Sprintf(":%d", SERVER_PORT), nil)
  if err != nil {
		panic(err)
	}
}
