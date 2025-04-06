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
	"github.com/VincNT21/kallaxy/client/models"
)

func (pm *GuiPageManager) GetLoginWindow() {
	// Create the window
	w := pm.appGui.NewWindow("Kallaxy Login")
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(800, 600))

	// Create UI objects
	// Texts
	pageTitleText := canvas.NewText("Please login", color.White)
	pageTitleText.Alignment = fyne.TextAlignCenter
	pageTitleText.TextSize = 40
	pageTitleText.TextStyle.Bold = true

	statusLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})

	// Entries
	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Username")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Pasword")

	// Buttons
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
				dialog.ShowError(err, w)
			}
		}

		if err == models.ErrUnauthorized {

		} else {
			log.Printf("--GUI-- User %v logged in\n", usernameEntry.Text)
			pm.GetHomeWindow()
			w.Close()
		}

	})

	passwordLostButton := widget.NewButtonWithIcon("Password lost", theme.QuestionIcon(), func() {})

	createNewUserButton := widget.NewButtonWithIcon("Create New User", theme.ContentAddIcon(), func() {
		log.Printf("--GUI-- User %v wants to create a new user\n", usernameEntry.Text)
		// Call the Create user window
		pm.GetCreateUserWindow(func() {
			// This part only runs if in Create user window, user confirm
			statusLabel.SetText("New user created, please login")
		})
	})

	exitButton := widget.NewButtonWithIcon("Exit App", theme.CancelIcon(), func() {
		dialog.ShowConfirm("Exit", "Are you sure you want to exit Kallaxy App ?", func(b bool) {
			if b {
				w.Close()
			}
		}, w)
	})

	// Group objects in VBox container
	objectsContainer := container.NewVBox(usernameEntry, passwordEntry, loginButton, passwordLostButton, createNewUserButton)

	// Create the global frame
	globalContainer := container.NewBorder(
		pageTitleText,
		exitButton,
		customSpacerHorizontal(100),
		customSpacerHorizontal(100),
		container.NewVBox(layout.NewSpacer(), objectsContainer, layout.NewSpacer()),
	)

	// Set container to window and show it
	w.SetContent(globalContainer)
	w.Show()
}

func (pm *GuiPageManager) GetBackWindow() {
	// Create the window
	w := pm.appGui.NewWindow("Kallaxy back")
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(800, 600))

	// Create objects
	// Text
	pageTitleText := canvas.NewText(fmt.Sprintf("Welcome back %s !", pm.appCtxt.APIClient.CurrentUser.Username), color.White)
	pageTitleText.Alignment = fyne.TextAlignCenter
	pageTitleText.TextSize = 40
	pageTitleText.TextStyle.Bold = true

	// Buttons
	enterButton := widget.NewButtonWithIcon("Enter app", theme.LoginIcon(), func() {
		// Call for RefreshTokens
		_, err := pm.appCtxt.APIClient.Auth.RefreshTokens()
		if err != nil {
			log.Println("--GUI-- Error with RefeshTokens")
			switch err {
			case models.ErrUnauthorized:
				dialog.ShowInformation("Error", "You need to login", w)
			case models.ErrServerIssue:
				dialog.ShowInformation("Error", "Error with server, please retry later", w)
			default:
				dialog.ShowError(err, w)
			}
		}
		pm.GetHomeWindow()
		w.Close()
	})
	notMeButton := widget.NewButtonWithIcon("Not you?", theme.CancelIcon(), func() {
		pm.GetLoginWindow()
		w.Close()
	})

	exitButton := widget.NewButtonWithIcon("Exit App", theme.CancelIcon(), func() {
		dialog.ShowConfirm("Exit", "Are you sure you want to exit Kallaxy App ?", func(b bool) {
			if b {
				w.Close()
			}
		}, w)
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

	// Set container to window and show it
	w.SetContent(globalContainer)
	w.Show()
}
