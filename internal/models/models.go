// Foodlebug models
package models

import (
	"time"
)

type User struct {
	Id       int
	Username string
	Password string
	Rating   float64
}

type Post struct {
	Id          int
	PosterId    int
	Photo       string
	Title       string
	Content     string
	TimePosted  time.Time
	LocationLat float64
	LocationLon float64
	Comments    []Comment
	Visible     bool
}

type Comment struct {
	Id       int
	PostId   int
	PosterId int
	Content  string
	Visible  bool
}

type Session struct {
	UserKey   string
	SessionId string
	CSRFToken string
	Expires   time.Time
	Created   time.Time
	IPAddress string
	UserAgent string
}
