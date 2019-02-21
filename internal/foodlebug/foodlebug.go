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

func TestHandler(w http.ResponseWriter, r *http.Request){
  w.WriteHeader(http.StatusOK)
  fmt.Fprintf(w, "hello")
}

func (f *Foodlebug) Run(){
  fmt.Println("Start...")

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

  users := f.store.GetUsers()

  for _, user := range users {
    fmt.Println(user.Username)
  }

  // routing
  r := mux.NewRouter()
  r.HandleFunc("/", TestHandler)
  http.Handle("/", r)

  // start server
  fmt.Printf("Starting server at %d...\n", SERVER_PORT)
  err := http.ListenAndServe(fmt.Sprintf(":%d", SERVER_PORT), nil)
  if err != nil {
		panic(err)
	}
}
