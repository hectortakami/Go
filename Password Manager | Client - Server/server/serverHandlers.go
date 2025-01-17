package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/A01377647/GoLang-Password-Manager/model"
	"github.com/A01377647/GoLang-Password-Manager/server/database"
	"github.com/A01377647/GoLang-Password-Manager/utils"
)

// función para escribir una respuesta del servidor
func response(w http.ResponseWriter, code int, payloadJSON string) {
	w.WriteHeader(code)
	fmt.Fprintf(w, payloadJSON)
}

// Añade un usuario a la BD
func registroUsuario(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	email := req.Form.Get("email")
	pass := req.Form.Get("pass")

	// Logs
	utils.AddLog("registroUsuario: [" + email + ", " + pass + "]")

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")

	// Añadimos el usuario a la base de datos
	if err := database.CreateUser(email, pass); err != nil {

		// Si ha ocurrido un error al añadir el usuario, comprobamos
		// el error y respondemos con el código http adecuado
		switch err.Error() {
		case "user already exists":
			response(w, 409, "") // (409 - Conflict)
		default:
			response(w, 500, "") // (500 - Internal Server Error)
		}

	} else {
		// Si la inserción se ha realizado correctamente
		response(w, 201, "")
	}
}

// Comprueba si existe un usuario en la BD
func loginUsuario(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	email := req.Form.Get("email")
	passw := req.Form.Get("pass")

	// Logs
	utils.AddLog("loginUsuario: [" + email + ", " + passw + "]")

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")

	// Añadimos el usuario a la base de datos
	if _, err := database.GetUser(email, passw); err != nil {

		// Si ha ocurrido un error al recuperar el usuario, comprobamos
		// el error y respondemos con el código http adecuado
		switch err.Error() {
		case "user not found":
			response(w, 404, "") // (404 - Not found)
		case "passwords do not match":
			response(w, 400, "") // (400 - Bad Request)
		default:
			response(w, 500, "") // (500 - Internal Server Error)
		}

	} else {
		// Si el usuario existe
		token, _ := CreateUserSession(email, false)
		response(w, 200, token)
	}
}

// Recupera las cuentas de servicio de un usuario de la BD
func listarEntradas(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	token := req.Form.Get("token")

	// Logs
	utils.AddLog("listarCuentas: [" + token + "]")

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")

	// Respondemos
	if email, errSession := GetUserFromSession(token); errSession != nil {
		// La sesión ha caducado o no es valida
		response(w, 401, "") // (401 - Unauthorized)
	} else if user, errUser := database.ReadUser(email); errUser != nil {
		response(w, 500, "") // (500 - Internal Server Error)
	} else {
		entriesList := model.ListaEntradas{}
		for entry := range user.Vault {
			tempEntry := user.Vault[entry]
			if tempEntry.Mode == 0 { // Texto
				// Guardamos solo lo que mostraremos, el título
				entriesList.Texts = append(entriesList.Texts, entry)
			} else if tempEntry.Mode == 1 { //Account
				// Guardamos solo lo que mostraremos, el título
				entriesList.Accounts = append(entriesList.Accounts, entry)
			} else if tempEntry.Mode == 2 { //Cards
				// Guardamos solo lo que mostraremos, el título
				entriesList.Cards = append(entriesList.Cards, entry)
			}
		}

		// Devolvemos la información
		if entriesJSON, errJSON := json.Marshal(entriesList); errJSON != nil {
			response(w, 500, "") // (500 - Internal Server Error)
		} else {
			response(w, 200, string(entriesJSON))
		}
	}
}

// Añade una entrada a un usuario de la BD
func crearEntrada(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	token := req.Form.Get("token")
	tituloEntrada := req.Form.Get("tituloEntrada")
	mode := req.Form.Get("mode") // Indica el tipo de entrada

	// Logs
	utils.AddLog("crearCuenta: [" + token + ", " + tituloEntrada + ", " + mode + "]")

	// Recogemos el email del usuario
	if email, errSession := GetUserFromSession(token); errSession != nil {
		// La sesión ha caducado o no es valida
		response(w, 401, "") // (401 - Unauthorized)
	} else {

		var errCreate error
		// Comprobamos el tipo de entrada que estamos creando
		if mode == "0" {
			// Si es una entrada de tipo texto
			textoEntrada := req.Form.Get("textoEntrada")
			errCreate = database.CreateTextVaultEntry(email, tituloEntrada, textoEntrada)

		} else if mode == "1" {
			// Si es una entrada de tipo cuenta de usuario
			usuarioEntradaCuenta := req.Form.Get("usuarioCuenta")
			passwordEntradaCuenta := req.Form.Get("passwordCuenta")
			errCreate = database.CreateAccountVaultEntry(email, tituloEntrada, usuarioEntradaCuenta, passwordEntradaCuenta)
		} else if mode == "2" {
			// Si es una entrada de tipo tarjeta de credito
			cardDigitsEntradaCuenta := req.Form.Get("cardDigits")
			cardExpirationEntradaCuenta := req.Form.Get("cardExpiration")
			cardCCVEntradaCuenta := req.Form.Get("cardCCV")
			errCreate = database.CreateCardVaultEntry(email, tituloEntrada, cardDigitsEntradaCuenta, cardExpirationEntradaCuenta, cardCCVEntradaCuenta)
		}

		// Respondemos
		if errCreate != nil {

			// Si ha ocurrido un error al insetar, comprobamos
			// el error y respondemos con el código http adecuado
			switch errCreate.Error() {
			case "user not found":
				response(w, 404, "") // (404 - Not found)
			case "entry already exists":
				response(w, 409, "") // (409 - Conflict)
			default:
				response(w, 500, "") // (500 - Internal Server Error)
			}

		} else {
			// Devolvemos la información
			response(w, 201, "")
		}
	}
}

// Recupera los detalles de una entrada concreta
func detallesEntrada(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	token := req.Form.Get("token")
	tituloEntrada := req.Form.Get("tituloEntrada")

	// Logs
	utils.AddLog("detallesEntrada: [" + token + ", " + tituloEntrada + "]")

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")

	// Respondemos
	if email, errSession := GetUserFromSession(token); errSession != nil {
		// La sesión ha caducado o no es valida
		response(w, 401, "") // (401 - Unauthorized)
	} else if entry, errRead := database.ReadVaultEntry(email, tituloEntrada); errRead != nil {

		// Si ha ocurrido un error al insetar, comprobamos
		// el error y respondemos con el código http adecuado
		switch errRead.Error() {
		case "user not found":
			response(w, 404, "") // (404 - Not found)
		case "entry not found":
			response(w, 404, "") // (404 - Not found)
		default:
			response(w, 500, "") // (500 - Internal Server Error)
		}

	} else {
		// Devolvemos la información
		if entryJSON, errJSON := json.Marshal(entry); errJSON != nil {
			response(w, 500, "") // (500 - Internal Server Error)
		} else {
			response(w, 200, string(entryJSON))
		}
	}
}

// Elimina una entrada de un usuario de la BD
func eliminarEntrada(w http.ResponseWriter, req *http.Request) {

	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	token := req.Form.Get("token")
	tituloEntrada := req.Form.Get("tituloEntrada")

	// Logs
	utils.AddLog("eliminarEntrada: [" + token + ", " + tituloEntrada + "]")

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")

	// Respondemos
	if email, errSession := GetUserFromSession(token); errSession != nil {
		// La sesión ha caducado o no es valida
		response(w, 401, "") // (401 - Unauthorized)
	} else if errDelete := database.DeleteVaultEntry(email, tituloEntrada); errDelete != nil {

		// Si ha ocurrido un error al borrar, comprobamos
		// el error y respondemos con el código http adecuado
		switch errDelete.Error() {
		case "user not found":
			response(w, 404, "") // (404 - Not found)
		case "entry not found":
			response(w, 404, "") // (404 - Not found)
		default:
			response(w, 500, "") // (500 - Internal Server Error)
		}

	} else {
		// Devolvemos la información
		response(w, 200, "")
	}
}
