package gui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/VincNT21/kallaxy/client/models"
)

func (pm *GuiPageManager) GetCreateUserWindow() {
	// Create the window
	w := pm.appGui.NewWindow("Create User")
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(500, 200))

	// Create objects
	titleLabel := widget.NewLabelWithStyle("Please provide information for user creation", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	statusLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})
	usernameEntry := widget.NewEntry()
	passwordEntry := widget.NewPasswordEntry()
	emailEntry := widget.NewEntry()

	// Group objects in a Form
	userForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Username", Widget: usernameEntry},
			{Text: "Password", Widget: passwordEntry},
			{Text: "Email", Widget: emailEntry},
		},
		OnSubmit: func() {
			log.Printf("--GUI-- CreateUser Form submitted - Username: %s, Email: %s\n", usernameEntry.Text, emailEntry.Text)
			dialog.ShowConfirm(
				"Confirmation",
				fmt.Sprintf("Are you sure to create this user ?\nUsername: %s\nPassword: %s\nEmail: %s", usernameEntry.Text, passwordEntry.Text, emailEntry.Text),
				func(b bool) {
					if b {
						_, err := pm.appCtxt.APIClient.Users.CreateUser(usernameEntry.Text, passwordEntry.Text, emailEntry.Text)
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
							w.Close()
						}

					}
				},
				w,
			)
		},
		OnCancel: func() {
			log.Println("--GUI-- CreateUser Form cancelled")
			w.Close()
		},
	}

	globalContainer := container.NewVBox(layout.NewSpacer(), titleLabel, layout.NewSpacer(), userForm, layout.NewSpacer(), statusLabel, layout.NewSpacer())

	w.SetContent(globalContainer)
	w.Show()
}
