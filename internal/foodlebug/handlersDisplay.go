// General display handlers
package foodlebug

import (
	"net/http"

	"github.com/jasmaa/foodlebug/internal/auth"
	"github.com/jasmaa/foodlebug/internal/models"
	"github.com/jasmaa/foodlebug/internal/store"
)

// Handle home page
func handleHome(store *store.Store) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		_, err := auth.SessionToUser(store, r)
		posts, _ := store.GetPosts()

		if err != nil {
			displayPage(w, "assets/templates/home.html", false, posts)
			return
		}

		displayPage(w, "assets/templates/home.html", true, posts)
	})
}

// Handle about page
func handleAbout(store *store.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := auth.SessionToUser(store, r)

		if err != nil {
			displayPage(w, "assets/templates/about.html", false, "")
			return
		}

		displayPage(w, "assets/templates/about.html", true, "")
	})
}

// Handle nearby page
func handleBrowse(store *store.Store) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		_, err := auth.SessionToUser(store, r)
		isLoggedIn := err == nil
		posts, _ := store.GetPosts()

		displayPage(w, "assets/templates/browse.html", isLoggedIn, posts)
	})
}

// Handle nearby page
func handleNearby(store *store.Store) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		_, err := auth.SessionToUser(store, r)
		isLoggedIn := err == nil

		//messages := make([]string, 10)
		posts, _ := store.GetPosts()
		//closePosts := make([]*models.Post, len(posts))

		displayPage(w, "assets/templates/nearby.html", isLoggedIn, posts)

		/*
			switch r.Method {
			case "GET":
				displayPage(w, "assets/templates/nearbyPrompt.html", isLoggedIn, messages)

			case "POST":
				lat, err := strconv.ParseFloat(r.FormValue("lat"), 64)
				if err != nil {
					messages := append(messages, "Could not submit post.")
					displayPage(w, "assets/templates/nearbyPrompt.html", isLoggedIn, messages)
					return
				}

				lon, err := strconv.ParseFloat(r.FormValue("lon"), 64)
				if err != nil {
					messages := append(messages, "Could not submit post.")
					displayPage(w, "assets/templates/nearbyPrompt.html", isLoggedIn, messages)
					return
				}

				dist, err := strconv.ParseFloat(r.FormValue("dist"), 64)
				if err != nil {
					messages := append(messages, "Could not submit post.")
					displayPage(w, "assets/templates/nearbyPrompt.html", isLoggedIn, messages)
					return
				}

				// filter out far away posts
				src := haversine.Coord{lat, lon}
				for _, p := range posts {
					target := haversine.Coord{p.LocationLat, p.LocationLon}
					mi, _ := haversine.Distance(src, target)

					if mi <= dist {
						closePosts = append(closePosts, p)
					}
				}

				displayPage(w, "assets/templates/nearby.html", isLoggedIn, closePosts)
			}
		*/
	})
}

// Handles user profile
func handleProfile(store *store.Store) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user, err := auth.SessionToUser(store, r)

		// Redirect to login on session failure
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// get posts made by user
		posts, _ := store.GetPosts()
		userPosts := make([]*models.Post, len(posts))
		for _, p := range posts {

			if p.PosterId == user.Id {
				userPosts = append(userPosts, p)
			}
		}

		type Data struct {
			User      *models.User
			UserPosts []*models.Post
		}
		displayPage(w, "assets/templates/profile.html", true, Data{user, userPosts})
	})
}
