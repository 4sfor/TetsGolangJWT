package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

// маршруты и сервер
func Server() {
	//маршрут получения токенов
	http.HandleFunc("/get_access", SetTokenHandler)
	//мартурт обновления токенов
	RefreshTokens := http.HandlerFunc(RefreshTokens)
	http.Handle("/refresh", RefreshTokens)

	log.Fatal(http.ListenAndServe(":8080", nil))

}

// функция отправки токенов
func SetTokenHandler(w http.ResponseWriter, r *http.Request) {
	//получаем токены
	tokenStringA, _ := TokenAccess()
	tokenStringR, _ := TokenRefresh()
	//записываем токен обновления в формате base64
	tokenStringR = base64.StdEncoding.EncodeToString([]byte(tokenStringR))
	//получаем guid пользователя из тела запроса
	guid := GetGUID(w, r)
	//ищем пользователя и передам токен обнолвения bcrypt
	searchUser := CheckGUID(guid, HashRefresh(tokenStringR))
	//отправляем токены в файлы cookie
	SetCookie(w, r, searchUser, tokenStringA, tokenStringR)
	return
}

// определяем файлы cookie
func SetCookie(w http.ResponseWriter, r *http.Request, searchUser bool, tokenStringA string, tokenStringR string) {
	if r.Method == "POST" {
		if searchUser == true {
			cookie := http.Cookie{
				Name:     "Access",
				Value:    tokenStringA,
				HttpOnly: true,
			}
			http.SetCookie(w, &cookie)
			cookie2 := http.Cookie{
				Name:     "Refresh",
				Value:    tokenStringR,
				HttpOnly: true,
			}
			http.SetCookie(w, &cookie2)
		} else {
			fmt.Fprintf(w, "пользователь не найден")
		}
	} else {
		fmt.Fprintf(w, "POST only")
	}

}

// записываем токен обновления как bcrypt
func HashRefresh(token string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(token), 14)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

// структура тела запроса
type User struct {
	GUID string `json:"guid"`
}

// получаем guid
func GetGUID(w http.ResponseWriter, r *http.Request) string {
	var user User
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&user)
		if err != nil {
			panic(err)
		}
		fmt.Println(user.GUID)
		return user.GUID

	}
	return ""
}

// получаем токен обновления из cookie
func GetRefreshToken(w http.ResponseWriter, r *http.Request) string {
	tokenRefreshCookie, err := r.Cookie("Refresh")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tokenRefreshCookie)
	return tokenRefreshCookie.Value

}

// обновляем токены
func RefreshTokens(w http.ResponseWriter, r *http.Request) {
	//вызываем функцию дял получения токена из cookie
	token := GetRefreshToken(w, r)
	//вызываем функцию для получения guid пользователя из тела запроса
	guid := GetGUID(w, r)
	//получаем токены
	tokenStringA, _ := TokenAccess()
	tokenStringR, _ := TokenRefresh()
	tokenStringR = base64.StdEncoding.EncodeToString([]byte(tokenStringR))
	//вызываем функцию для проверки того, что переданный в cookie токен обнлвения соответсвует тому что храниться в базе
	searchToken := CheckRefresh(guid, token)
	if r.Method == "POST" {
		if searchToken == true {
			//устанавливаем cookie
			cookie := http.Cookie{
				Name:     "Access",
				Value:    tokenStringA,
				HttpOnly: true,
			}
			http.SetCookie(w, &cookie)
			cookie2 := http.Cookie{
				Name:     "Refresh",
				Value:    tokenStringR,
				HttpOnly: true,
			}
			http.SetCookie(w, &cookie2)
			//сохраняем новый токен  обновления в базу
			CheckGUID(guid, HashRefresh(tokenStringR))
		} else {
			fmt.Fprintf(w, "токен не найден")
		}
	} else {
		fmt.Fprintf(w, "POST only")
	}
}
