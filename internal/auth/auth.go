// Foodlebug authentication
package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/jasmaa/foodlebug/internal/models"
	"github.com/jasmaa/foodlebug/internal/store"
)

// Generate n random bits
func GenerateRandomBits(n int) string {
	result := make([]byte, n/8)
	_, err := io.ReadFull(rand.Reader, result)
	if err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(result)
}

// Create new session for a user
func NewSession(store *store.Store, user *models.User, expires time.Time, ipAddress, userAgent string) (*models.Session, error) {
	if expires.IsZero() {
		expires = time.Now().AddDate(0, 0, 3)
	}
	s := &models.Session{
		UserKey:   user.Username,
		SessionId: GenerateRandomBits(128),
		CSRFToken: GenerateRandomBits(256),
		Expires:   expires,
		Created:   time.Now(),
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}

	// Update or create
	var err error

	err = store.DeleteSessions(s.UserKey)
	if err != nil {
		return nil, err
	}

	err = store.InsertSession(s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func SessionToUser(store *store.Store, r *http.Request) (*models.User, error) {

	cookies := r.Cookies()
	cookieValue := ""

	// Read in cookie value
	for i := range cookies {
		if cookies[i].Name == "foodlebug" {
			if cookies[i].Value != "" {
				cookieValue = cookies[i].Value
				break
			}
		}
	}

	if cookieValue == "" {
		return nil, errors.New("no cookie")
	}

	var userKey string
	var sessionId string

	vals := strings.Split(cookieValue, "_")
	userKey = vals[0]
	sessionId = vals[1]

	s, err := store.GetSession(userKey)

	if s.SessionId != sessionId {
		return nil, errors.New("session invalid")
	}

	user, err := store.GetUser("username", userKey)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Hash password
func HashPassword(password string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Verify password hash
func VerifyPassword(hash, plain string) bool {

	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(plain),
	)
	if err != nil {
		return false
	}
	return true
}

// Create new user
func CreateNewUser(store *store.Store, username string, password string) error {

	// hash password
	hash, err := HashPassword(password)
	if err != nil {
		return err
	}

	// add user to db
	err = store.AddUser(&models.User{
		Id:       store.GenerateUserId(),
		Username: username,
		Password: hash,
		Rating:   0,
	})

	return err
}
