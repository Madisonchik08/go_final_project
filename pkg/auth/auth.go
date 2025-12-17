package auth

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// SigninRequest represents login request payload.
// SigninRequest представляет данные запроса на вход.
type SigninRequest struct {
	Password string `json:"password"`
}

// SigninResponse represents login response payload.
// SigninResponse представляет данные ответа на вход.
type SigninResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

// SigninHandler handles user authentication request.
// SigninHandler обрабатывает запрос на аутентификацию пользователя.
func SigninHandler(w http.ResponseWriter, r *http.Request) {
	var req SigninRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, SigninResponse{Error: "Неверный формат"}, http.StatusBadRequest)
		return
	}

	expectedPass := os.Getenv("TODO_PASSWORD")
	if expectedPass == "" {
		// If password is not set, authentication is not required
		// Если пароль не задан, аутентификация не требуется
		respondJSON(w, SigninResponse{Error: "Аутентификация не требуется"}, http.StatusBadRequest)
		return
	}

	if req.Password != expectedPass {
		respondJSON(w, SigninResponse{Error: "Неверный пароль"}, http.StatusUnauthorized)
		return
	}

	token, err := GenerateJWT(expectedPass)
	if err != nil {
		respondJSON(w, SigninResponse{Error: "Ошибка генерации токена"}, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		MaxAge:   8 * 3600,
		HttpOnly: true,
	})

	respondJSON(w, SigninResponse{Token: token}, http.StatusOK)
}

// GenerateJWT creates a JWT token for authenticated user.
// GenerateJWT создает JWT токен для аутентифицированного пользователя.
func GenerateJWT(password string) (string, error) {
	secret := os.Getenv("TODO_JWT_SECRET")
	if secret == "" {
		secret = "my_secret_key"
	}

	claims := &jwt.MapClaims{
		"password_hash": HashPassword(password),
		"exp":           time.Now().Add(8 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// HashPassword creates a simple hash from password.
// HashPassword создает простой хэш из пароля.
func HashPassword(p string) string {
	return string(p[0]) + string(p[len(p)-1]) + string(rune(len(p)))
}

// respondJSON writes JSON response with specified status code.
// respondJSON записывает JSON ответ с указанным кодом статуса.
func respondJSON(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
	}
}
