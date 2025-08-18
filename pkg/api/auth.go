package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// claims - это кастомная структура полезных данных (payload) для нашего JWT
// Она встраивает стандартные RegisteredClaims (для полей вроде "exp" - срок годности и т.д.)
// и добавляет наше собственное поле PasswordHash для сверки с текущим паролем
type claims struct {
	jwt.RegisteredClaims
	PasswordHash string `json:"password_hash"`
}

// getPasswordHash вычисляет хэш SHA256 для строки пароля
// В данном проекте это используется как "контрольная сумма" пароля, чтобы не хранить сам пароль в токене.
// Но мы также имеем возможность проверить, что токен был выдан для текущего активного пароля
func getPasswordHash(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

// generateJWT создаёт новый JWT-токен, подписанный паролем.
// evnPass - текущий пароль из переменной окружения, который используется
// для подписи, и для генерации хэша в payload.
func generateJWT(envPass string) (string, error) {
	// Создаём payload для токена.
	claims := claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		},	
		PasswordHash: getPasswordHash(envPass),
	}

	// Создаём новый токен с нашими полезными данными (payload) и методом подписи HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Подписываем токен. В качестве секрета используется текущий пароль
	tokenString, err := token.SignedString([]byte(envPass))
	if err != nil {
		return "", fmt.Errorf("Ошибка при генерации токена: %w", err)
	}

	return tokenString, nil
}

// auth - это middleware (обёртка для обработчика) для проверки аутентификации.
// Он "оборачивает" другой обработчик (next) и выполняет проверку токена
// перед тем, как передать управление основному обработчику
func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		envPass := os.Getenv("TODO_PASSWORD")

		// если пароль не установлен, аутентификая не требуется
		if envPass == "" {
			next(w, r)
			return
		}

		// Получаем cookie с именем "token"
		cookie, err := r.Cookie("token")
		// Если cookie нет, значит пользователь не аутентифицирован
		if err != nil {
			http.Error(w, "Требуется аутентификация", http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value
		claims := &claims{}

		// Парсим токен.
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			// проверка на правильный метод подписи
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("неожиданный метод подписи: %v", t.Header["alg"])
			}

			// возвращаем текущий пароль в качестве секрета для сверки подписи
			return []byte(envPass), nil

		})

		if err != nil {
			http.Error(w, "Невалидный токен", http.StatusUnauthorized)
			return
		}

		// Доп проверка, что токен помечен как валидный
		if !token.Valid {
			http.Error(w, "Невалидный токен", http.StatusUnauthorized)
			return
		}

		// Сверяем хэш пароля из токена с хэшем текущего пароля и переменной окружения
		// Если пароль на сервере изменился, эта проверка провалится
		// и старый токен будет недействительным, несмотря на правильную подпись и не истёкший срок
		if claims.PasswordHash != getPasswordHash(envPass) {
			http.Error(w, "Невалидный токен (пароль был изменен)", http.StatusUnauthorized)
			return
		}

		// Все проверки пройдены. Вызываем обработчик
		next(w, r)
	}
}

// loginHandler обрабатывает POST-запросы на /api/signin для входа в систему.
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJson(w, http.StatusMethodNotAllowed, map[string]string{"error": "Неправильный метод HTTP-запроса"})
		return
	}

	// Декодируем JSON с паролем из тела запроса
	var reqData map[string]string
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Ошибка десериализации JSON"})
		return
	}

	receivedPassword, ok := reqData["password"]
	if !ok {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "поле 'password' отсутствует"})
		return
	}

	envPass := os.Getenv("TODO_PASSWORD")
	// Защитная проверка
	if envPass == "" {
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": "аутентификая не настроена на сервере"})
		return
	}

	// Сверяем полученный пароль с текущим из переменной окружения
	if receivedPassword != envPass {
		writeJson(w, http.StatusUnauthorized, map[string]string{"error": "Неверный пароль"})
		return
	}

	// Если пароли совпали, генерируем JWT
	tokenString, err := generateJWT(envPass)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	cookie := &http.Cookie{
		Name: "token",
		Value: tokenString,
		Path: "/",
		Expires: time.Now().Add(8 * time.Hour),
	}

	// Устанавливаем токен в http-only cookie
	http.SetCookie(w, cookie)

	// Возвращаем токен в теле JSON-ответа
	writeJson(w, http.StatusOK, map[string]string{"token": tokenString})
}
