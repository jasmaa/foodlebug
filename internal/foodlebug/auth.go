package foodlebug

import(
  "fmt"
  "net/http"
  "html/template"
  "golang.org/x/crypto/bcrypt"

  "github.com/jasmaa/foodlebug/internal/store"
  "github.com/jasmaa/foodlebug/internal/models"
)

func hashPassword(password string) (string, error) {
  // Hash password
  hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
  if err != nil {
    return "", err
  }
  return string(hash), nil
}

func verifyPassword(hash, plain string) bool {
  // Verify password hash
  err := bcrypt.CompareHashAndPassword(
    []byte(hash),
    []byte(plain),
  )
  if err != nil {
    return false
  }
  return true
}

func createNewUser(store *store.Store, username string, password string) error{
  // New user
  // hash password
  hash, err := hashPassword(password)
  if err != nil {
    return err
  }

  // add user to db
  err = store.AddUser(&models.User{
    Id:store.GenerateUserId(),
    Username:username,
    Password:hash,
    Rating:0,
  })

  return err
}

func displayCreateAccount(w http.ResponseWriter, messages []string) {
  // Display create account screen get
  w.WriteHeader(http.StatusOK)
  var t *template.Template
  t, _ = template.ParseFiles(
    "assets/templates/main.html",
    "assets/templates/footer.html",
    "assets/templates/createAccount.html",
  )
  t.ExecuteTemplate(w, "main", messages)
}

func displayLogin(w http.ResponseWriter, messages []string) {
  // Display login screen get
  w.WriteHeader(http.StatusOK)
  var t *template.Template
  t, _ = template.ParseFiles(
    "assets/templates/main.html",
    "assets/templates/footer.html",
    "assets/templates/login.html",
  )
  t.ExecuteTemplate(w, "main", messages)
}

func displayCreateAccountSuccess(w http.ResponseWriter) {
  // Account creation success
  w.WriteHeader(http.StatusOK)
  var t *template.Template
  t, _ = template.ParseFiles(
    "assets/templates/main.html",
    "assets/templates/footer.html",
    "assets/templates/createAccountSuccess.html",
  )
  t.ExecuteTemplate(w, "main", "")
}

// === HANDLERS ===

func handleCreateAccount(store *store.Store) http.Handler {
  // Account creation page
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

    messages := make([]string, 10)

    switch r.Method {
      case "GET":
        displayCreateAccount(w, messages)

      case "POST":
        username := r.FormValue("username")
        password := r.FormValue("password")
        confirm := r.FormValue("confirm-password")

        if username == "" || password == "" {
          messages := append(messages, "Fields cannot be empty.")
          displayCreateAccount(w, messages)
          return
        } else if password != confirm {
          messages := append(messages, "Passwords do not match.")
          displayCreateAccount(w, messages)
          return
        }

        err := createNewUser(store, username, password)

        if err != nil {
          messages := append(messages, "Account could not be created.")
          displayCreateAccount(w, messages)
          return
        }

        // go to success page
        displayCreateAccountSuccess(w)
    }
  })
}

func handleLogin(store *store.Store) http.Handler {
  // Login page
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

    messages := make([]string, 10)

    switch r.Method {
      case "GET":
        displayLogin(w, messages)

      case "POST":
        username := r.FormValue("username")
        password := r.FormValue("password")

        user, err := store.GetUser("username", username)
        if err != nil {
          messages := append(messages, "User does not exist.")
          displayLogin(w, messages)
          return
        }

        if !verifyPassword(user.Password, password) {
          messages := append(messages, "Incorrect password.")
          displayLogin(w, messages)
          return
        }

        // do session stuff here!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
        fmt.Fprintf(w, "%s\n%s\n%s", user.Username, user.Password, password)
    }
  })
}
