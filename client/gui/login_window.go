package gui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func getLoginWindow(a fyne.App) {

	// Creates the window
	w := a.NewWindow("Welcome")
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(300, 200))

	// Creates objects
	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Username")
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Pasword")
	loginButton := widget.NewButtonWithIcon("Login", theme.ConfirmIcon(), func() {
		log.Printf("User %v wants to login\n", usernameEntry.Text)
	})
	createNewUserButton := widget.NewButtonWithIcon("Create New User", theme.ContentAddIcon(), func() {
		log.Printf("User %v wants to create a new user\n", usernameEntry.Text)
	})

	// Groups objects in VBox container
	objectsContainer := container.NewVBox(usernameEntry, passwordEntry, loginButton, createNewUserButton)

	// Create the center row
	centerRow := container.NewHBox(layout.NewSpacer(), objectsContainer, layout.NewSpacer())

	// Create the global frame
	welcomeLabel := widget.NewLabel("Welcome to Project Kallaxy App !")
	globalContainer := container.NewVBox(layout.NewSpacer(), welcomeLabel, centerRow, layout.NewSpacer())

	// Set container to window
	w.SetContent(globalContainer)
	w.Show()
}
