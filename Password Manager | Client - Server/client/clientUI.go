package client

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/A01377647/GoLang-Password-Manager/model"
	"github.com/A01377647/GoLang-Password-Manager/utils"
	"github.com/fatih/color"
)

var httpClient *http.Client

func startUI(c *http.Client) {
	httpClient = c
	printWelcomeMenu()
	uiMainMenu("", "")
}

// Pantalla de bienvenida con las opciones de
// login, registro y cerrar aplicación
func uiMainMenu(showError string, showSuccess string) {

	// Limpiamos la pantalla
	utils.ClearScreen()

	// Mensaje de confirmación de acción en caso de existir
	if showSuccess != "" {
		color.HiGreen("* %s\n\n", showSuccess)
	}
	// Opciones
	fmt.Printf(" \nPlease, select an option: \n\n")
	fmt.Println("\t[1] Login User")
	fmt.Println("\t[2] Register")
	fmt.Println("\t[q] Quit")

	// Mensaje de error en caso de existir
	if showError != "" {
		color.HiRed("\n[Error] %s", showError)
	}

	// Lectura de opción elegida
	color.HiBlue(" \n[Option] ")
	inputSelectionStr := utils.CustomScanf()

	// Ejecución de la opción elegida
	switch {
	case inputSelectionStr == "1": // Login
		uiLoginUser("")
	case inputSelectionStr == "2": // Registro
		uiRegisterUser("")
	case inputSelectionStr == "q", inputSelectionStr == "quit": // Salir
		printByeMsg()
		os.Exit(0)
	default:
		uiMainMenu("Invalid option!", "")
	}
}

// Pantalla de creación de usuario
func uiRegisterUser(showError string) {

	// Limpiamos la pantalla
	utils.ClearScreen()

	// Título de la pantalla
	color.HiBlue("\n\n* Registering new user\n\n")

	// Mensaje de error en caso de existir
	if showError != "" {
		color.HiRed("[Error Register] %s\n\n", showError)
	}

	// Lectura de datos del nuevo usuario
	fmt.Print("Please, provide an email for your account\n")
	color.HiBlue("[Email] ")
	inputUser := utils.CustomScanf()
	fmt.Print("\nPlease, provide a master password to access all your records and data saved\n")
	color.HiBlue("[Master Password] ")
	inputPass := utils.CustomScanf()

	// Petición al servidor
	if err := registroUsuario(httpClient, inputUser, inputPass); err != nil {
		// yes hay un error, mostramos el mensaje de error adecuado
		switch err.Error() {
		case "user already exists":
			uiRegisterUser("This email already exists in database.")
		default:
			uiRegisterUser("An error ocurred while the new user was created. Please try again!")
		}
	} else {
		// Registro completado, volvemos a la pantalla de inicio
		uiMainMenu("", "New user correctly created. Please login to enter session\n")
	}
}

// Pantalla de entrada de usuarios
func uiLoginUser(showError string) {

	// Limpiamos la pantalla
	utils.ClearScreen()

	// Título de la pantalla
	color.HiBlue("\n\n* Login as existing user\n\n")

	// Mensaje de error en caso de existir
	if showError != "" {
		color.HiRed("* %s\n\n", showError)
	}

	// Lectura de datos del usuario
	fmt.Print("Please, enter your user credentials\n")
	color.HiBlue("[Email] ")
	inputUser := utils.CustomScanf()
	color.HiBlue("[Master Password] ")
	inputPass := utils.CustomScanf()

	// Petición al servidor
	if err := loginUsuario(httpClient, inputUser, inputPass); err != nil {

		// yes hay un error, mostramos el mensaje de error adecuado
		switch err.Error() {
		case "user not found":
			uiMainMenu("Invalid user, no register found with the provided data. Try again!", "")
		case "passwords do not match":
			// No damos información detallada del error en este caso
			uiMainMenu("Invalid user, no register found with the provided data. Try again!", "")
		default:
			uiMainMenu("Sorry! Something went wrong while login :(", "")
		}

	} else {
		// Login completado, vamos a la pantalla principal del usuario
		uiUserMainMenu("", "")
	}
}

// Pantalla principal del usuario, listado de entradas
func uiUserMainMenu(showError string, showSuccess string) {

	// Limpiamos la pantalla
	utils.ClearScreen()

	// Mensaje de confirmación de acción en caso de existir
	if showSuccess != "" {
		color.HiGreen("\n* %s\n", showSuccess)
	}

	// Recuperamos las cuentas del usuario

	color.HiBlue("\n*********** KEY-SET DATA VAULT ***********\n")
	color.HiBlue("******************************************\n\n")
	// Petición al servidor
	entradas, err := listarEntradas(httpClient)
	if err != nil {
		// yes hay un error, mostramos el mensaje de error adecuado
		switch err.Error() {
		case "unauthorized":
			uiLoginUser("Your session has expired.")
		default:
			fmt.Println("Something went wrong while returning your data vault :(" + err.Error())
		}
	} else {
		// Mostramos la lista de cuentas de usuario guardadas

		boldBlue := color.New(color.FgHiCyan, color.Bold)
		if (entradas.Accounts != nil && len(entradas.Accounts) != 0) ||
			(entradas.Texts != nil && len(entradas.Texts) != 0) ||
			(entradas.Cards != nil && len(entradas.Cards) != 0) {

			// Mostramos la lista de cuentas de usuario guardadas
			if entradas.Accounts != nil && len(entradas.Accounts) != 0 {
				boldBlue.Printf(" Accounts/Passwords\n")
				// Imprimimos los resultados
				for c := range entradas.Accounts {
					color.HiYellow("    [%s]\n", entradas.Accounts[c])
				}
				fmt.Printf("\n")
			}

			// Mostramos la lista de tarjetas guardadas
			if entradas.Cards != nil && len(entradas.Cards) != 0 {
				boldBlue.Printf(" Credit/Debit Cards\n")
				// Imprimimos los resultados
				for c := range entradas.Cards {
					color.HiYellow("    [%s]\n", entradas.Cards[c])
				}
				fmt.Printf("\n")
			}

			// Mostramos la lista de textos guardados
			if entradas.Texts != nil && len(entradas.Texts) != 0 {
				boldBlue.Printf(" Notes\n")
				// Imprimimos los resultados
				for c := range entradas.Texts {
					color.HiYellow("    [%s]\n", entradas.Texts[c])
				}
			}

		} else {
			color.HiCyan("          [No records store yet]\n")
		}
	}
	color.HiBlue("\n******************************************\n")
	color.HiBlue("*********** KEY-SET DATA VAULT ***********\n\n")

	// Opciones
	fmt.Println(" [1] Add a new Key-Set to Vault")
	fmt.Println(" [2] View my Key-Set data")
	fmt.Println(" [0] Log Out")

	// Mensaje de error en caso de existir
	if showError != "" {
		color.HiRed("\n* %s", showError)
	}

	// Lectura de opción elegida
	color.HiBlue("\n[Option] ")
	inputSelectionStr := utils.CustomScanf()

	// Ejecución de la opción elegida
	switch {
	case inputSelectionStr == "1":
		uiAddNewEntry("")
	case inputSelectionStr == "2":
		fmt.Print("Provide the register name\n ")
		color.HiBlue("[Key-Set ID]")
		inputEntrySelectionStr := utils.CustomScanf()
		uiDetailsEntry("", inputEntrySelectionStr)
	case inputSelectionStr == "0":
		uiMainMenu("", "")
	default:
		uiUserMainMenu("Invalid option! Try again", "")
	}
}

// Pantalla de creación de nueva entrada
func uiAddNewEntry(showError string) {

	// Limpiamos la pantalla
	utils.ClearScreen()

	// Título de la pantalla
	fmt.Printf("# Please, select an option...\n\n")

	// Solicitamos información de lo que queremos guardar de entre las poyesbles
	fmt.Println("[1] Password Account")
	fmt.Println("[2] Credit/Debit Card Number")
	fmt.Println("[3] Text Note (Raw Text)")

	fmt.Println("[0] Return Main Menu")

	// Mensaje de error en caso de existir
	if showError != "" {
		color.HiRed("\n* %s\n", showError)
	}

	// Lectura de opción elegida
	color.HiBlue("\n[Option] ")
	inputEntryMode := utils.CustomScanf()

	switch inputEntryMode {
	case "1":
		uiAddNewAccountEntry("")
	case "2":
		uiAddNewCardEntry("")
	case "3":
		uiAddNewTextEntry("")
	case "0":
		uiUserMainMenu("", "")
	default:
		uiAddNewEntry("Invalid option selected.")
	}
}

// Pantalla de creación de nueva entrada de tipo texto
func uiAddNewTextEntry(showError string) {

	// Limpiamos la pantalla
	utils.ClearScreen()

	// Título de la pantalla
	color.HiBlue("* Add new key-set [text note] to vault\n\n")

	// Mensaje de error en caso de existir
	if showError != "" {
		color.HiRed("[ERROR NOTE] %s\n\n", showError)
	}

	// Lectura de los datos de la nueva entrada
	fmt.Printf("\nEnter a title for your note: ")
	inputTitle := utils.CustomScanf()
	fmt.Printf("\nWrite the body content of your note and then hit [ENTER] to store:\n\n")
	inputText := utils.CustomScanf()

	// Petición al servidor
	if err := crearEntradaDeTexto(httpClient, inputTitle, inputText); err != nil {
		// yes hay un error, mostramos el mensaje de error adecuado
		switch err.Error() {
		case "unauthorized":
			uiLoginUser("The session expired.")
		case "user not found":
			uiLoginUser("An error occurred while saving the note.")
		case "entry already exists":
			uiUserMainMenu("An existing note already have this title. ", "")
		default:
			uiUserMainMenu("An error occured while saving your note :(", "")
		}
	} else {
		uiUserMainMenu("", "The note '["+inputTitle+"]' was successfully stored :) ")
	}
}

func uiAddNewCardEntry(showError string) {

	// Limpiamos la pantalla
	utils.ClearScreen()

	// Título de la pantalla
	color.HiBlue("* Add new key-set [credit/debit card] to vault\n\n")

	// Mensaje de error en caso de existir
	if showError != "" {
		color.HiRed("[ERROR CARD-SET] %s\n\n", showError)
	}

	// Lectura de los datos de la nueva entrada
	fmt.Print("Card description (HSBC, BBVA, 'My lovely card'): ")
	inputCardType := utils.CustomScanf()
	color.HiBlue("[Card 16 digits] (format XXXX-XXXX-XXXX-XXXX) ")
	inputCardDigits := utils.CustomScanf()
	color.HiBlue("[Expiration date] (format MM/YY)")
	inputCardExpiration := utils.CustomScanf()
	color.HiBlue("[CCV] (format 123)")
	inputCardCCV := utils.CustomScanf()

	// Petición al servidor
	if err := createCardSet(httpClient, inputCardType, inputCardDigits, inputCardExpiration, inputCardCCV); err != nil {
		// yes hay un error, mostramos el mensaje de error adecuado
		switch err.Error() {
		case "unauthorized":
			uiLoginUser("The session expired.")
		case "user not found":
			uiLoginUser("An error occurred while saving the card credentials.")
		case "entry already exists":
			uiUserMainMenu("An existing card credentials already have this title.", "")
		default:
			uiUserMainMenu("An error occured while saving your card credentials :(", "")
		}
	} else {
		uiUserMainMenu("", "The ["+inputCardType+"] credentials were stored correctly :)")
	}
}

// Pantalla de creación de nueva entrada de tipo cuenta de usuario
func uiAddNewAccountEntry(showError string) {

	// Limpiamos la pantalla
	utils.ClearScreen()

	// Título de la pantalla
	color.HiBlue("* Add new key-set [password account] to vault\n\n")

	// Mensaje de error en caso de existir
	if showError != "" {
		color.HiRed("[ERROR ACCOUNT-SET] %s\n\n", showError)
	}

	// Lectura de los datos de la nueva entrada
	fmt.Print("Account description (facebook, instagram, twitter...): ")
	inputAccountType := utils.CustomScanf()
	color.HiBlue("[Account User] ")
	inputAccountUser := utils.CustomScanf()
	fmt.Print("Do you want the system to provide you a random password? (yes, no): ")
	inputGeneratePassw := utils.CustomScanf()

	var finalPassw string
	if inputGeneratePassw == "yes" || inputGeneratePassw == "y" {

		// Solicitamos información de como se desea generar la contraseña
		for {
			// Tamaño de la contraseña
			var genLenght int
			for {
				fmt.Print("Define the password length (no. of chars): ")
				inputLenght := utils.CustomScanf()
				if convLenght, err := strconv.Atoi(inputLenght); err == nil {
					genLenght = convLenght
					break
				}
			}

			// La contraseña generada puede tener números
			fmt.Print("Do you want to include numbers (Alphanumeric)? (yes, no): ")
			inputWithNums := utils.CustomScanf()
			genWithNums := inputWithNums == "yes" || inputWithNums == "y"

			// La contraseña generada puede tener yesmbolos
			fmt.Print("Do you want to include symbols? (yes, no): ")
			inputWithSymbols := utils.CustomScanf()
			genWithSymbols := inputWithSymbols == "yes" || inputWithSymbols == "y"

			// Mostramos la contraseña y preguntamos al usuario yes está de acuerdo
			finalPassw = utils.GeneratePassword(genLenght, true, genWithNums, genWithSymbols)
			color.HiBlue("Your new auto-generated password will be [ %s ]\nAre you agree? (yes, no): ", finalPassw)
			inputConfirm := utils.CustomScanf()
			if inputConfirm == "yes" || inputConfirm == "y" {
				break
			}
		}

	} else {
		color.HiBlue("[Password]: ")
		finalPassw = utils.CustomScanf()
	}

	// Petición al servidor
	if err := crearEntradaDeCuenta(httpClient, inputAccountType, inputAccountUser, finalPassw); err != nil {
		// yes hay un error, mostramos el mensaje de error adecuado
		switch err.Error() {
		case "unauthorized":
			uiLoginUser("The session expired.")
		case "user not found":
			uiLoginUser("An error occurred while saving the account credentials.")
		case "entry already exists":
			uiUserMainMenu("An existing account-set already have this title.", "")
		default:
			uiUserMainMenu("An error occured while saving your account-set :(", "")
		}
	} else {
		uiUserMainMenu("", "The ["+inputAccountType+"] credentials were stored correctly :)")
	}
}

// Pantalla de visualización de detalles de una entrada
func uiDetailsEntry(showError string, entryName string) {

	// Limpiamos la pantalla
	utils.ClearScreen()

	// Título de la pantalla
	color.HiGreen("\n\nKey-Set [ %s ] details\n\n", entryName)

	// Petición al servidor
	color.HiCyan("****************************************\n\n")
	entry, err := detallesEntrada(httpClient, entryName)
	if err != nil {
		// yes hay un error, mostramos el mensaje de error adecuado
		switch err.Error() {
		case "unauthorized":
			uiLoginUser("The session expired.")
		case "not found":
			uiUserMainMenu("There was a problem retrieving key-set details", "")
		default:
			fmt.Println("There was a problem retrieven key-set" + err.Error())
		}
	} else {
		// yes los detalles de la cuenta están vacios
		if (model.VaultEntry{}) == entry {
			// Volvemos al menú del usuario
			uiUserMainMenu("There was a problem retrieving key-set details from this user account.", "")
		}

		// Comprobamos el tipo de entrada (texto, cuenta) y la mostramos
		if entry.Mode == 0 {
			// yes es una entrada de tipo texto
			color.HiYellow("[Note Content] \n\n%s\n", entry.Text)

		} else if entry.Mode == 1 {
			// yes es una entrada de tipo cuenta de usuario
			color.HiYellow("[User] -> %s \n", entry.User)
			color.HiYellow("[Password] -> %s \n", entry.Password)

		} else if entry.Mode == 2 {
			// yes es una entrada de tipo tarjeta
			color.HiYellow("[Card Number] -> %s \n", entry.CardDigits)
			color.HiYellow("[Expiration Date] -> %s \n", entry.CardExpiration)
			color.HiYellow("[CCV] -> %s \n", entry.CardCCV)
		}
	}
	color.HiCyan("\n****************************************\n\n")

	// Opciones
	fmt.Println("[1] Delete")
	fmt.Println("[0] Return")

	// Mensaje de error en caso de existir
	if showError != "" {
		color.HiRed("\n[ERROR DELETE KEY-SET] %s", showError)
	}

	// Lectura de opción elegida
	color.HiBlue("\n[Option] ")
	inputSelectionStr := utils.CustomScanf()

	switch {
	case inputSelectionStr == "1":
		fmt.Print("Do you really want to delete this key-set register? This action is permanent and irreversible (yes, no): ")
		inputDecisyeson := utils.CustomScanf()
		if inputDecisyeson == "yes" || inputDecisyeson == "y" {

			// Petición al servidor para eliminar la entrada de la BD
			if errDel := eliminarEntrada(httpClient, entryName); errDel != nil {
				// yes hay un error, mostramos el mensaje de error adecuado
				switch errDel.Error() {
				case "unauthorized":
					uiLoginUser("The session expired.")
				case "not found":
					uiUserMainMenu("Error while deleting register.", "")
				default:
					fmt.Println("Error while deleting register." + err.Error())
				}

			} else {
				// Se ha eliminado correctamente
				uiUserMainMenu("", "Key-set ["+entryName+"] removed successfully")
			}
		} else {
			uiDetailsEntry("", entryName)
		}
	case inputSelectionStr == "0":
		uiUserMainMenu("", "")
	default:
		uiDetailsEntry("Invalid option", entryName)
	}
}

func printWelcomeMenu() {

	msg := "\n\n*************************************************************************************************\n" +
		"*************************************************************************************************\n" +
		"    \t\t    ____                                                         __       \n" +
		"   \t\t   / __ \\  ____ _   _____   _____ _      __  ____    _____  ____/ /       \n" +
		"  \t\t  / /_/ / / __ `/  / ___/  / ___/| | /| / / / __ \\  / ___/ / __  /        \n" +
		" \t\t / ____/ / /_/ /  (__  )  (__  ) | |/ |/ / / /_/ / / /    / /_/ /         \n" +
		"\t\t/_/ __  _\\_,_|   /____/  /____/  |__/|__/  \\____/ /_/     \\__,_/          \n" +
		"   \t\t   /  |/  /  ____ _   ____   ____ _   ____ _  ___    _____                \n" +
		"  \t\t  / /|_/ /  / __ `/  / __ \\ / __ `/  / __ `/ / _ \\  / ___/                \n" +
		" \t\t / /  / /  / /_/ /  / / / // /_/ /  / /_/ / /  __/ / /                    \n" +
		"\t\t/_/  /_/   \\__,_|  /_/ /_/ \\__,_|   \\__, /  \\___/ /_/                     \n" +
		"\t\t                                   /____/                                 \n" +
		"*************************************************************************************************\n" +
		"*************************************************************************************************\n"

	color.HiBlue(msg)

}

func printByeMsg() {
	bye := "\n\n*************************************************************************************************\n" +
		"*************************************************************************************************\n" +
		"	   \t\t\t    __                     __       \n" +
		"	   \t\t\t   / /_    __  __  ___    / /       \n" +
		"	  \t\t\t  / __ \\  / / / / / _ \\  / /      \n" +
		"	 \t\t\t / /_/ / / /_/ / /  __/ /_/         \n" +
		"	\t\t\t/_.___/  \\__, /  \\___/ (_)        \n" +
		"		\t\t\t/____/                      \n\n" +
		"*************************************************************************************************\n" +
		"*************************************************************************************************\n"
	color.HiBlue(bye)
}
