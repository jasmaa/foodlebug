// Database control
package store

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"

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

// Find a free id
func (store *Store) GenerateId(table string) int {

	res := true
	var id int32

	// keep looking for id
	for res {

		id = rand.Int31()

		rows, err := store.db.Query("SELECT exists (SELECT 1 FROM "+table+" WHERE id=$1 LIMIT 1)", id)
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

// === USER ===

// Insert user into db
func (store *Store) InsertUser(u *models.User) error {

	var err error
	_, err = store.db.Query(
		"INSERT INTO users (id, username, password, rating) VALUES ($1, $2, $3, $4)",
		u.Id, u.Username, u.Password, u.Rating)
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
			return nil, err
		}
		return user, nil
	}

	return nil, errors.New("Could not get user")
}

// Retrieves all users
func (store *Store) GetUsers() ([]*models.User, error) {

	// Retrieve user
	rows, err := store.db.Query("SELECT id, username, password, rating FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*models.User{}

	// Iterate thru users
	for rows.Next() {
		user := &models.User{}
		err = rows.Scan(&user.Id, &user.Username, &user.Password, &user.Rating)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// === SESSION ===

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

// Remove all sessions with userkey
func (store *Store) DeleteSessions(userKey string) error {
	_, err := store.db.Query(fmt.Sprintf("DELETE FROM sessions WHERE userKey='%s'", userKey))
	if err != nil {
		return err
	}
	return nil
}

// Retrieve session by userkey
func (store *Store) GetSession(userKey string) (*models.Session, error) {

	rows, err := store.db.Query(fmt.Sprintf("SELECT sessionId, CSRFToken, expires, created, ipAddress, userAgent FROM sessions WHERE userKey='%s'", userKey))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expireRaw string
	var createdRaw string

	if rows.Next() {
		s := &models.Session{UserKey: userKey}
		err = rows.Scan(
			&s.SessionId,
			&s.CSRFToken,
			&expireRaw,
			&createdRaw,
			&s.IPAddress,
			&s.UserAgent,
		)
		if err != nil {
			return nil, err
		}

		s.Expires, _ = time.Parse(time.RFC3339, expireRaw)
		s.Created, _ = time.Parse(time.RFC3339, createdRaw)

		return s, nil
	}

	return nil, err
}

// Update session
func (store *Store) UpdateSessionExpire(userKey string, expires time.Time) error {
	rows, err := store.db.Query("UPDATE sessions SET expires=$1 WHERE userKey=$2", expires, userKey)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

// === POST ===

// Insert post
func (store *Store) InsertPost(p *models.Post) error {
	var err error
	_, err = store.db.Query(
		"INSERT INTO posts (id, posterId, photo, title, content, timePosted, locationLat, locationLon, visible) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		p.Id, p.PosterId, p.Photo, p.Title, p.Content, p.TimePosted, p.LocationLat, p.LocationLon, p.Visible)
	if err != nil {
		return err
	}

	return nil
}

// Retrieves all posts in reverse chronological order
func (store *Store) GetPosts() ([]*models.Post, error) {

	// Retrieve user
	rows, err := store.db.Query("SELECT id, posterId, photo, title, content, timePosted, locationLat, locationLon, visible FROM posts ORDER BY timePosted DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*models.Post{}

	// Iterate thru users
	for rows.Next() {
		post := &models.Post{}
		err = rows.Scan(
			&post.Id,
			&post.PosterId,
			&post.Photo,
			&post.Title,
			&post.Content,
			&post.TimePosted,
			&post.LocationLat,
			&post.LocationLon,
			&post.Visible,
		)
		if err != nil {
			return nil, err
		}

		if post.Visible {
			posts = append(posts, post)
		}
	}

	return posts, nil
}

// Get post by field
func (store *Store) GetPost(field string, val string) (*models.Post, error) {

	rows, err := store.db.Query(fmt.Sprintf("SELECT id, posterId, photo, title, content, timePosted, locationLat, locationLon, visible FROM posts WHERE %s='%s'", field, val))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		post := &models.Post{}
		err = rows.Scan(
			&post.Id,
			&post.PosterId,
			&post.Photo,
			&post.Title,
			&post.Content,
			&post.TimePosted,
			&post.LocationLat,
			&post.LocationLon,
			&post.Visible,
		)
		if err != nil {
			return nil, err
		}
		return post, nil
	}

	return nil, err
}
