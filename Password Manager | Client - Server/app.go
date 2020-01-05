package main

import (
	"fmt"
	"os"

	"github.com/A01377647/GoLang-Password-Manager/client"
	"github.com/A01377647/GoLang-Password-Manager/server"
	"github.com/A01377647/GoLang-Password-Manager/utils"
)

func main() {

	// Recogemos el valor de los argumentos
	if len(os.Args) == 2 {

		argMode := os.Args[1]
		if argMode == "client" {
			client.Start()
		} else if argMode == "server" {
			server.Launch()
		} else {
			fmt.Printf("Invalid argument mode! Try again\n")
		}

	} else if (len(os.Args)) == 4 {

		argMode := os.Args[1]
		if argMode == "logger" {
			argInput := os.Args[2]
			argOutput := os.Args[3]
			utils.LaunchLogger(argInput, argOutput)
		} else {
			fmt.Printf("Invalid argument mode! Try again\n")
		}
	} else {
		fmt.Printf("Invalid parameters detected! Try again\n")
	}
}
