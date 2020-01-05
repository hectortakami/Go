package config

// AppName contiene el nombre de la aplicación
var AppName = "Password Manager"

// SecureServerPort Puerto Seguro Cliente
var SecureServerPort = ":10443"

// SecureURL Url segura
var SecureURL = "https://127.0.0.1"

// MaxTimeSession es el tiempo máximo de sesión (segundos)
var MaxTimeSession = 60 * 30

// PassDBEncrypt es la clave de cifrado del fichero de base de datos
var PassDBEncrypt = []byte("RXR0stzTuxSq1DszjPVkwRU0l5d8KYZs")

// EncryptLogs se encarga de indicar si se desea cifrar el log del servidor
var EncryptLogs = false

// PassEncryptLogs Clave de cifrado de los ficheros de logs
var PassEncryptLogs = []byte("RXR0stzTuxSq1DszjPVkwRU0l5d8KYZs")
