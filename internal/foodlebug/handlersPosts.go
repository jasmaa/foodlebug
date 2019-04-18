// Handlers for food posts
package foodlebug

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jasmaa/foodlebug/internal/auth"
	"github.com/jasmaa/foodlebug/internal/models"
	"github.com/jasmaa/foodlebug/internal/store"
)

// Handle post page
func handlePage(store *store.Store) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		post, err := store.GetPost("id", vars["postId"])
		if err != nil || post == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// check if logged in
		_, err = auth.SessionToUser(store, r)
		if err != nil {
			displayPage(w, "assets/templates/postPage.html", false, post)
			return
		}

		displayPage(w, "assets/templates/postPage.html", true, post)
	})
}

// Post entry page handler
func handlePostEntry(store *store.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user, err := auth.SessionToUser(store, r)

		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		messages := make([]string, 10)

		switch r.Method {
		case "GET":
			displayPage(w, "assets/templates/postEntry.html", true, messages)

		case "POST":
			entryTitle := r.FormValue("entry-title")
			content := r.FormValue("entry-content")
			entryPhoto := r.FormValue("entry-photo")

			lat, err := strconv.ParseFloat(r.FormValue("lat"), 64)
			if err != nil {
				messages := append(messages, "Could not submit post.")
				displayPage(w, "assets/templates/postEntry.html", false, messages)
				return
			}

			lon, err := strconv.ParseFloat(r.FormValue("lon"), 64)
			if err != nil {
				messages := append(messages, "Could not submit post.")
				displayPage(w, "assets/templates/postEntry.html", false, messages)
				return
			}

			// Check if entries filled
			if entryTitle == "" || content == "" || entryPhoto == "" {
				messages := append(messages, "Could not submit post.")
				displayPage(w, "assets/templates/postEntry.html", false, messages)
				return
			}

			// Get b64 data
			data := strings.Split(entryPhoto, ",")
			if len(data) != 2 {
				messages := append(messages, "Could not submit post.")
				displayPage(w, "assets/templates/postEntry.html", false, messages)
				return
			}

			// Limit size
			if float64(len(data[1]))/1.37 > 5000000 {
				messages := append(messages, "Image exceeded 500KB.")
				displayPage(w, "assets/templates/postEntry.html", false, messages)
				return
			}

			// Insert post into db
			post := &models.Post{
				Id:          store.GenerateId("posts"),
				PosterId:    user.Id,
				Photo:       data[1],
				Title:       entryTitle,
				Content:     content,
				TimePosted:  time.Now(),
				LocationLat: lat,
				LocationLon: lon,
				Visible:     true,
			}

			err = store.InsertPost(post)
			if err != nil {
				messages := append(messages, "Could not submit post.")
				displayPage(w, "assets/templates/postEntry.html", false, messages)
				return
			}

			// Redirect to home
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	})
}
