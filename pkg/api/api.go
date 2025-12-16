package api

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"go_final_project/pkg/auth"
)

// DateLayout is the expected date format for parameters.
const DateLayout = "20060102"

// authMiddleware проверяет аутентификацию пользователя.
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwt string
			cookie, err := r.Cookie("token")
			if err == nil {
				jwt = cookie.Value
			}
			var valid bool

			// Пароль задан, проверяем токен
			if jwt != "" {
				// Парсим JWT токен
				token, err := jwt.ParseWithClaims(jwt, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
					secret := os.Getenv("TODO_JWT_SECRET")
					if secret == "" {
						secret = "my_secret_key"
					}
					return []byte(secret), nil
				})

				if err == nil && token.Valid {
					if claims, ok := token.Claims.(*jwt.MapClaims); ok {
						passwordHash, ok := (*claims)["password_hash"].(string)
						if ok {
							// Сравниваем хэш пароля из токена с текущим хэшем пароля
							currentHash := auth.HashPassword(pass)
							valid = (passwordHash == currentHash)
						}
					}
				}
			}

			if !valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}

// Init registers all API handlers on provided mux.
func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", nextDateHandler)
	mux.HandleFunc("/api/signin", auth.SigninHandler)
	mux.HandleFunc("/api/task", authMiddleware(taskHandler))
	mux.HandleFunc("/api/tasks", authMiddleware(tasksHandler))
	mux.HandleFunc("/api/task/done", authMiddleware(doneTaskHandler))
}
