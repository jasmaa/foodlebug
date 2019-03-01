// Database control

package store

import (
  "fmt"
  "errors"
  "math/rand"
  "database/sql"
  _ "github.com/lib/pq"

  "github.com/jasmaa/foodlebug/internal/models"
)

type Store struct {
  db *sql.DB
}

func (store *Store) Connect(host string, port int, user string, password string, dbname string) {
  // Inits db connection
  // connect to db
  var err error
  connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
  store.db, err = sql.Open("postgres", connString)
  if err != nil {
    panic(err)
  }

  // ping db to see if actually connected
  err = store.db.Ping()
  if err != nil {
    panic(err)
  }
}

func (store *Store) AddUser(u *models.User) error {
  // Insert user into db
  var err error
  _, err = store.db.Query(
    "INSERT INTO users (id, username, password, rating) VALUES ($1, $2, $3, $4)",
    u.Id, u.Username, u.Password, u.Rating
  )
  if err != nil {
    return err
  }

  return nil
}

func (store *Store) GetUser(field string, val string) (*models.User, error) {
  // Retrieve user by field
  rows, err := store.db.Query(fmt.Sprintf("SELECT id, username, password, rating FROM users WHERE %s='%s'", field, val))
	if err != nil {
    return nil, errors.New("Could not query db")
	}
	defer rows.Close()

  if rows.Next(){
    user := &models.User{}
    err = rows.Scan(&user.Id, &user.Username, &user.Password, &user.Rating)
    if err != nil {
      panic(err)
    }
    return user, nil;
  }

  return nil, errors.New("Could not get user")
}

func (store *Store) GetUsers() []*models.User {
  // Retrieve user
  rows, err := store.db.Query("SELECT id, username, password, rating FROM users")
	if err != nil {
    panic(err)
	}
	defer rows.Close()

  users := []*models.User{}

  // Iterate thru users
  for rows.Next() {
    user := &models.User{}
    err = rows.Scan(&user.Id, &user.Username, &user.Password, &user.Rating)
    if err != nil {
      panic(err)
    }

    users = append(users, user)
  }

  return users
}

func (store *Store) GenerateUserId() int {
  // Find a free id
  res := true
  var id int32

  // keep looking for id
  for res {

    id = rand.Int31()

    rows, err := store.db.Query("SELECT exists (SELECT 1 FROM users WHERE id=$1 LIMIT 1)", id)
  	if err != nil {
      panic(err)
  	}
  	defer rows.Close()

    rows.Next()
    err = rows.Scan(&res)
    if err != nil{
      panic(err)
    }
  }

  return int(id);
}
