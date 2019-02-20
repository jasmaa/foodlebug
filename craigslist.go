package main

import (
  "fmt"
  "database/sql"

  _ "github.com/lib/pq"
)

const (
  host = "localhost"
  port = 5432
  user  = "postgres"
  password = "<put passwd here>"
  dbname = "craigslist_db"
)

func main(){
  // connect to db
  connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
  db, err := sql.Open("postgres", connString)
  if err != nil {
    panic(err)
  }
  defer db.Close()

  // ping db to see if actually connected
  err = db.Ping()
  if err != nil {
    panic(err)
  }

  db.Query("SELECT * from craigslist_db")
}
