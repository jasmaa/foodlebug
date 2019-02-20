// Models

package main

import (
  "time"
)

type User struct {
  Id int
  Username string
  Password string
  Rating float64
}

type Post struct {
  Id int
  PosterId int
  Photo []byte
  Content string
  TimePosted time.Duration
  LocationLat float64
  LocationLon float64
  Comments []Comment
  Visible bool
}

type Comment struct {
  Id int
  PostId int
  PosterId int
  Content string
  Visible bool
}
