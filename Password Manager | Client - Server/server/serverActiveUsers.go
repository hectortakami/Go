package server

import (
	"errors"
	"time"

	"github.com/A01377647/GoLang-Password-Manager/config"
	"github.com/A01377647/GoLang-Password-Manager/model"
	"github.com/A01377647/GoLang-Password-Manager/utils"
)

var activeUsers = make(map[string]*model.ActiveUser)

// CreateUserSession añade un usuario a la lista de usuarios
// activos asignandole un token de sesión generado.
func CreateUserSession(userEmail string, useA2F bool) (string, string) {

	cleanAllInactiveUsers()

	var token = generateSessionToken()
	useA2F = false

	activeUsers[token] = &model.ActiveUser{
		UserEmail:          userEmail,
		SesssionExpireTime: time.Now().Add(time.Second * time.Duration(config.MaxTimeSession)),
	}

	return token, ""
}

// GetUserFromSession recupera el correo electrónico del usuario
// si está activo a partir del token de sesión que se indica.
func GetUserFromSession(token string) (string, error) {

	var userEmail = ""
	var err error
	if tempUser, ok := activeUsers[token]; !ok {
		err = errors.New("session not found")
	} else if isSessionExpired(token) {
		delete(activeUsers, token)
		err = errors.New("session expired")
	} else {
		userEmail = tempUser.UserEmail
		resetSessionExpireTime(token)
	}

	return userEmail, err
}

func resetSessionExpireTime(token string) {
	if tempUser, ok := activeUsers[token]; ok {
		tempUser.SesssionExpireTime = time.Now().Add(time.Second * time.Duration(config.MaxTimeSession))
	}
}

// cleanInactiveUsers recorre la lista de usuarios activos y
// elimina aquellos cuya sesión haya cadudado.
func cleanAllInactiveUsers() {
	for k := range activeUsers {
		if isSessionExpired(k) {
			delete(activeUsers, k)
		}
	}
}

// isSessionExpired comprueba si el tiempo de validez de una
// sesión sigue estando activo
func isSessionExpired(token string) bool {
	isExpired := false
	if tempUser, ok := activeUsers[token]; ok {
		currentTime := time.Now()
		if currentTime.After(tempUser.SesssionExpireTime) {
			isExpired = true
		}
	} else {
		isExpired = true
	}
	return isExpired
}

// generateSessionToken genera un token para la sesión que será
// el que use el usuario al realizar las peticiones.
func generateSessionToken() string {
	tokenRaw, _ := utils.GenerateRandomBytes(24)
	tokenSrc := utils.EncodeBase64(tokenRaw)
	return tokenSrc
}
