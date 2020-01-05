package model

import "time"

/* ----------- DATABASE ----------- */

type Usuario struct {
	UserPassword     string
	UserPasswordSalt string
	Vault            map[string]VaultEntry
}

type VaultEntry struct {
	Mode int
	// Mode 0 - Plain text	// Mode 1 - Account // Mode 2 - Credit/Debit Card
	Text           string
	User           string
	Password       string
	CardDigits     string
	CardExpiration string
	CardCCV        string
}

/* Demo estructura en json (sin cifrados)
   "alu@alu.ua.es": {
       "UserPassword": "accoutPass",
       "UserPasswordSalt": "accoutSalrPass",
       "Vault": {
           "memoria": {
               "Mode": "0",
               "Text": "texto de la entrada"
           },
           "twitter": {
               "Mode": "1",
               "User": "usuarioTwitter",
               "Password": "54321"
           }
       }
   }
*/

/* -------------------------------- */

/*  ----- USUARIO ACTIVO ----- */
type ActiveUser struct {
	UserEmail          string
	SesssionExpireTime time.Time
}

/*  -------------------------- */

/*  ----- PETICIONES ----- */

type DetallesUsuario struct {
	Email      string
	NumEntries int
}

type ListaEntradas struct {
	Texts    []string
	Accounts []string
	Cards    []string
}

/* ----------------------- */
