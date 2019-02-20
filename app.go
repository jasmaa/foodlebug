// Main program

package main

import (
  "fmt"
)

const (
  host = "localhost"
  port = 5432
  user  = "postgres"
  password = "<password>"
  dbname = "craigslist_db"
)

func main(){
  fmt.Println("Start...")

  store := Store{}
  store.Connect(host, port, user, password, dbname)

  /*
  store.AddUser(&User{
    2,
    "John",
    "super secret password",
    6.9,
  })
  */

  users := store.GetUsers()

  for _, user := range users {
    fmt.Println(user.Username)
  }
}
