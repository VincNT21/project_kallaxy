package gui

import (
	"fmt"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/VincNT21/kallaxy/client/context"
	"github.com/VincNT21/kallaxy/client/models"
)

func createLoginContent(appCtxt *context.AppContext) *fyne.Container {
	// Create UI objects
	// Top title
	pageTitleText := canvas.NewText("Please login", color.White)
	pageTitleText.Alignment = fyne.TextAlignCenter
	pageTitleText.TextSize = 40
	pageTitleText.TextStyle.Bold = true

	// Bottom exit button
	exitButton := widget.NewButtonWithIcon("Exit App", theme.CancelIcon(), func() {
		dialog.ShowConfirm("Exit", "Are you sure you want to exit Kallaxy App ?", func(b bool) {
			if b {
				appCtxt.MainWindow.Close()
			}
		}, appCtxt.MainWindow)
	})

	// Status label
	statusLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})

	// Entries
	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Username")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Pasword")

	// Buttons
	loginButton := widget.NewButtonWithIcon("Login", theme.ConfirmIcon(), func() {
		buttonFuncLogin(appCtxt, usernameEntry, passwordEntry, statusLabel)
	})

	passwordLostButton := widget.NewButtonWithIcon("Password lost", theme.QuestionIcon(), func() {
		showPasswordLostSecondaryWindow(appCtxt)
	})

	createNewUserButton := widget.NewButtonWithIcon("Create New User", theme.ContentAddIcon(), func() {
		log.Printf("--GUI-- User %v wants to create a new user\n", usernameEntry.Text)
		// Call the Create user window
		appCtxt.PageManager.ShowCreateUserPage()
	})

	// Group objects in VBox container
	objectsContainer := container.NewVBox(
		usernameEntry,
		passwordEntry,
		loginButton,
		passwordLostButton,
		createNewUserButton,
		statusLabel,
	)

	// Create the global frame
	globalContainer := container.NewBorder(
		pageTitleText,
		exitButton,
		customSpacerHorizontal(100),
		customSpacerHorizontal(100),
		container.NewVBox(layout.NewSpacer(), objectsContainer, layout.NewSpacer()),
	)

	// Send content container back to page manager
	return globalContainer
}

func buttonFuncLogin(appCtxt *context.AppContext, usernameEntry, passwordEntry *widget.Entry, statusLabel *widget.Label) {
	_, err := appCtxt.APIClient.Auth.LoginUser(usernameEntry.Text, passwordEntry.Text)
	if err != nil {
		log.Printf("--GUI-- User %v failed to login\n", usernameEntry.Text)
		switch err {
		case models.ErrUnauthorized:
			statusLabel.SetText("Bad username/password")
			statusLabel.Refresh()
		case models.ErrServerIssue:
			statusLabel.SetText("Error with server, please retry later")
			statusLabel.Refresh()
		default:
			dialog.ShowError(err, appCtxt.MainWindow)
		}
	} else {
		log.Printf("--GUI-- User %v logged in\n", usernameEntry.Text)
		appCtxt.PageManager.ShowHomePage()
	}
}

func createWelcomeBackContent(appCtxt *context.AppContext) *fyne.Container {
	// Create objects
	// Top title
	pageTitleText := canvas.NewText(fmt.Sprintf("Welcome back %s !", appCtxt.APIClient.CurrentUser.Username), color.White)
	pageTitleText.Alignment = fyne.TextAlignCenter
	pageTitleText.TextSize = 40
	pageTitleText.TextStyle.Bold = true

	// Bottom exit button
	exitButton := widget.NewButtonWithIcon("Exit App", theme.CancelIcon(), func() {
		dialog.ShowConfirm("Exit", "Are you sure you want to exit Kallaxy App ?", func(b bool) {
			if b {
				appCtxt.MainWindow.Close()
			}
		}, appCtxt.MainWindow)
	})

	// Buttons
	enterButton := widget.NewButtonWithIcon("Enter app", theme.LoginIcon(), func() {
		// Call for RefreshTokens
		_, err := appCtxt.APIClient.Auth.RefreshTokens()
		if err != nil {
			log.Println("--GUI-- Error with RefeshTokens")
			dialog.ShowConfirm("Error", "There is a problem with server\nYou need to login again", func(b bool) {
				if b {
					appCtxt.PageManager.ShowLoginPage()
				}
			}, appCtxt.MainWindow)
		} else {
			appCtxt.PageManager.ShowHomePage()
		}
	})
	notMeButton := widget.NewButtonWithIcon("Not you?", theme.CancelIcon(), func() {
		appCtxt.PageManager.ShowLoginPage()
	})

	// Group Buttons
	buttonContainer := container.NewVBox(enterButton, notMeButton)

	// Create the global frame
	globalContainer := container.NewBorder(
		pageTitleText,
		exitButton,
		customSpacerHorizontal(100),
		customSpacerHorizontal(100),
		container.NewVBox(layout.NewSpacer(), buttonContainer, layout.NewSpacer()),
	)

	// Send content container back to page manager
	return globalContainer
}
