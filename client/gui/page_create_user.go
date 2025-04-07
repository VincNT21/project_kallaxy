package gui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/VincNT21/kallaxy/client/context"
	"github.com/VincNT21/kallaxy/client/models"
)

func createUserContent(appCtxt *context.AppContext) *fyne.Container {
	// Create objects
	// Texts
	titleLabel := widget.NewLabelWithStyle("Please provide information for user creation", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	statusLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})

	// Entries
	usernameEntry := widget.NewEntry()
	passwordEntry := widget.NewPasswordEntry()
	emailEntry := widget.NewEntry()

	// Group entries in a Form
	userForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Username", Widget: usernameEntry},
			{Text: "Password", Widget: passwordEntry},
			{Text: "Email", Widget: emailEntry},
		},
		// Forms buttons
		OnSubmit: func() {
			log.Printf("--GUI-- CreateUser Form submitted - Username: %s, Email: %s\n", usernameEntry.Text, emailEntry.Text)
			dialog.ShowConfirm(
				"Confirmation",
				fmt.Sprintf("Are you sure to create this user ?\nUsername: %s\nPassword: %s\nEmail: %s", usernameEntry.Text, passwordEntry.Text, emailEntry.Text),
				func(b bool) {
					if b {
						// If user confims, call CreateUser API function
						_, err := appCtxt.APIClient.Users.CreateUser(usernameEntry.Text, passwordEntry.Text, emailEntry.Text)
						if err != nil {
							switch err {
							case models.ErrBadRequest:
								statusLabel.SetText("User creation failed: one field was not provided")
							case models.ErrConflict:
								statusLabel.SetText("User creation failed: username or email already used")
							case models.ErrServerIssue:
								statusLabel.SetText("User creation failed: error with server, please retry later")
							default:
								statusLabel.SetText("User creation failed: unknown error")
							}
						} else {
							// Activate the callback onConfirm function
							dialog.ShowInformation("Created", "New user created !\nPlease login", appCtxt.MainWindow)
							appCtxt.PageManager.ShowLoginPage()
						}
					}
				},
				appCtxt.MainWindow,
			)
		},
		OnCancel: func() {
			log.Println("--GUI-- CreateUser Form cancelled")
			appCtxt.PageManager.ShowLoginPage()
		},
	}

	// Set the global frame
	globalContainer := container.NewVBox(layout.NewSpacer(), titleLabel, layout.NewSpacer(), userForm, layout.NewSpacer(), statusLabel, layout.NewSpacer())

	return globalContainer
}
