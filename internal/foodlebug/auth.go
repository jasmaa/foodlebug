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

func handleCreateAccount(store *store.Store) http.Handler {
  // Account creation page
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    switch r.Method {
      case "GET":
        w.WriteHeader(http.StatusOK)
        var t *template.Template
        t, _ = template.ParseFiles(
          "assets/templates/main.html",
          "assets/templates/footer.html",
          "assets/templates/createAccount.html",
        )
        t.ExecuteTemplate(w, "main", "")

      case "POST":
        username := r.FormValue("username")
        password := r.FormValue("password")

        // hash password
        hash, err := hashPassword(password)
        if err != nil {
          fmt.Fprintf(w, "An error has occured")
          return
        }

        // add user to db
        err = store.AddUser(&models.User{
          Id:store.GenerateUserId(),
          Username:username,
          Password:hash,
          Rating:0,
        })
        if err != nil {
          fmt.Fprintf(w, "An error has occured")
          return
        }

        // go to success page
        w.WriteHeader(http.StatusOK)
        var t *template.Template
        t, _ = template.ParseFiles(
          "assets/templates/main.html",
          "assets/templates/footer.html",
          "assets/templates/createAccountSuccess.html",
        )
        t.ExecuteTemplate(w, "main", "")
    }
  })
}

func handleLogin(store *store.Store) http.Handler {
  // Login page
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    switch r.Method {
      case "GET":
        w.WriteHeader(http.StatusOK)
        var t *template.Template
        t, _ = template.ParseFiles(
          "assets/templates/main.html",
          "assets/templates/footer.html",
          "assets/templates/login.html",
        )
        t.ExecuteTemplate(w, "main", "")

      case "POST":
        username := r.FormValue("username")
        password := r.FormValue("password")

        fmt.Fprintf(w, "%s\n%s", username, password)
    }
  })
}
