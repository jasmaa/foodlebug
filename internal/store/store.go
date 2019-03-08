// Database control
package store

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"

	_ "github.com/lib/pq"

	"github.com/jasmaa/foodlebug/internal/models"
)

type Store struct {
	db *sql.DB
}

// Inits db connection
func (store *Store) Connect(host string, port int, user string, password string, dbname string) {

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

// Insert user into db
func (store *Store) AddUser(u *models.User) error {

	var err error
	_, err = store.db.Query(
		"INSERT INTO users (id, username, password, rating) VALUES ($1, $2, $3, $4)",
		u.Id, u.Username, u.Password, u.Rating)
	if err != nil {
		return err
	}

	return nil
}

// Insert session into db
func (store *Store) InsertSession(s *models.Session) error {

	var err error
	_, err = store.db.Query(
		"INSERT INTO sessions (userKey, sessionId, CSRFToken, expires, created, ipAddress, userAgent) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		s.UserKey, s.SessionId, s.CSRFToken, s.Expires, s.Created, s.IPAddress, s.UserAgent)
	if err != nil {
		return err
	}

	return nil
}

// Retrieve user by field
func (store *Store) GetUser(field string, val string) (*models.User, error) {

	rows, err := store.db.Query(fmt.Sprintf("SELECT id, username, password, rating FROM users WHERE %s='%s'", field, val))
	if err != nil {
		return nil, errors.New("Could not query db")
	}
	defer rows.Close()

	if rows.Next() {
		user := &models.User{}
		err = rows.Scan(&user.Id, &user.Username, &user.Password, &user.Rating)
		if err != nil {
			panic(err)
		}
		return user, nil
	}

	return nil, errors.New("Could not get user")
}

// Retrieves all users
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

// Find a free id
func (store *Store) GenerateUserId() int {

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
		if err != nil {
			panic(err)
		}
	}

	return int(id)
}
