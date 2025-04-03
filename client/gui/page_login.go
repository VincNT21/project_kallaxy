package gui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/VincNT21/kallaxy/client/models"
)

func (pm *GuiPageManager) GetLoginWindow() {
	// Create the window
	w := pm.appGui.NewWindow("Welcome")
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(400, 300))

	// Create objects
	titleLabel := widget.NewLabelWithStyle("Welcome to Project Kallaxy App !", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	statusLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})

	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Username")
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Pasword")
	loginButton := widget.NewButtonWithIcon("Login", theme.ConfirmIcon(), func() {
		_, err := pm.appCtxt.APIClient.Auth.LoginUser(usernameEntry.Text, passwordEntry.Text)
		if err != nil {
			log.Printf("--GUI-- User %v failed to login\n", usernameEntry.Text)
			switch err {
			case models.ErrUnauthorized:
				statusLabel.SetText("Bad username/password")
			case models.ErrServerIssue:
				statusLabel.SetText("Error with server, please retry later")
			default:
				statusLabel.SetText("Couldn't login: unknown error")
			}
		}

		if err == models.ErrUnauthorized {

		} else {
			log.Printf("--GUI-- User %v logged in\n", usernameEntry.Text)
			pm.GetHomeWindow()
			w.Close()
		}

	})
	createNewUserButton := widget.NewButtonWithIcon("Create New User", theme.ContentAddIcon(), func() {
		log.Printf("--GUI-- User %v wants to create a new user\n", usernameEntry.Text)
		pm.GetCreateUserWindow()
		statusLabel.SetText("New user created, please login")
	})

	// Group objects in VBox container
	objectsContainer := container.NewVBox(usernameEntry, passwordEntry, loginButton, createNewUserButton)

	// Create the center row
	centerRow := container.NewHBox(layout.NewSpacer(), objectsContainer, layout.NewSpacer())

	// Create the global frame
	globalContainer := container.NewVBox(layout.NewSpacer(), titleLabel, layout.NewSpacer(), centerRow, layout.NewSpacer(), statusLabel, layout.NewSpacer())

	// Set container to window
	w.SetContent(globalContainer)
	w.Show()
}

func (pm *GuiPageManager) GetBackWindow() {
	// Create the window
	w := pm.appGui.NewWindow("Welcome back")
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(400, 300))

	// Create objects
	titleLabel := widget.NewLabelWithStyle(fmt.Sprintf("Welcome back %s !", pm.appCtxt.APIClient.LastUser.Username), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	enterButton := widget.NewButtonWithIcon("Enter app", theme.HomeIcon(), func() {
		pm.GetHomeWindow()
		w.Close()
	})
	notMeButton := widget.NewButtonWithIcon("Not you?", theme.CancelIcon(), func() {
		pm.GetLoginWindow()
		w.Close()
	})

	// Group Buttons
	buttonRow := container.NewHBox(layout.NewSpacer(), enterButton, notMeButton, layout.NewSpacer())

	// Create the global frame

	globalContainer := container.NewVBox(layout.NewSpacer(), titleLabel, layout.NewSpacer(), buttonRow, layout.NewSpacer())

	// Set container to window
	w.SetContent(globalContainer)
	w.Show()
}
