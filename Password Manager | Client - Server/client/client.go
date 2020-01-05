package client

import (
	"crypto/tls"
	"net/http"

	"github.com/A01377647/GoLang-Password-Manager/config"
)

var baseURL = config.SecureURL + config.SecureServerPort

// Start Inicio del cliente
func Start() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// Cliente http utilizado en la aplicación
	client := &http.Client{Transport: tr}

	// Lanzamiento de la interfaz
	startUI(client)
}
