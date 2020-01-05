package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/A01377647/GoLang-Password-Manager/config"
	"github.com/A01377647/GoLang-Password-Manager/server/database"
	"github.com/A01377647/GoLang-Password-Manager/utils"
	"github.com/fatih/color"
)

// Launch lanza el servidor
func Launch() {

	// suscripción SIGINT
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	// Rutas disponibles
	mux := http.NewServeMux()
	mux.Handle("/usuario/login", http.HandlerFunc(loginUsuario))
	mux.Handle("/usuario/registro", http.HandlerFunc(registroUsuario))
	mux.Handle("/vault", http.HandlerFunc(listarEntradas))
	mux.Handle("/vault/nueva", http.HandlerFunc(crearEntrada))
	mux.Handle("/vault/detalles", http.HandlerFunc(detallesEntrada))
	mux.Handle("/vault/eliminar", http.HandlerFunc(eliminarEntrada))

	srv := &http.Server{Addr: config.SecureServerPort, Handler: mux}

	color.HiGreen("System server online and listening :)")
	go func() {
		if err := srv.ListenAndServeTLS("server/certs/cert.pem", "server/certs/key.pem"); err != nil {

			log.Printf("[CONNECTION] %s\n", err)
		}
	}()

	<-stopChan // espera señal SIGINT
	color.HiYellow("Shutting down server ...")

	// Guarda la información de la BD en un fichero
	database.After()

	//Guarda logs en fichero
	utils.AfterLogs()

	// Apaga el servidor de forma segura
	ctx, fnc := context.WithTimeout(context.Background(), 5*time.Second)
	fnc()
	srv.Shutdown(ctx)

	color.HiGreen("Server went to sleep! Bye bye!")
}
