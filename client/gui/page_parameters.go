package gui

import (
	"fmt"
	"image/color"

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

func createParametersContent(appCtxt *context.AppContext) *fyne.Container {
	// Create UI objects
	// Texts
	statusLabel := widget.NewLabelWithStyle("You'll be automatically logged out if you update your personal info", fyne.TextAlignCenter, fyne.TextStyle{})
	titleText := canvas.NewText("Personal Informations", color.White)
	titleText.TextSize = 20
	titleText.Alignment = fyne.TextAlignCenter
	titleText.TextStyle.Bold = true

	usernameLabel := widget.NewLabel(fmt.Sprintf("Username: %s", appCtxt.APIClient.CurrentUser.Username))
	emailLabel := widget.NewLabel(fmt.Sprintf("Email: %s", appCtxt.APIClient.CurrentUser.Email))

	// Versions text
	clientVersion := widget.NewLabel(fmt.Sprintf("Client current version: %s", appCtxt.APIClient.ClientVersion))
	serverVersion := widget.NewLabel(fmt.Sprintf("Server current version: %s", appCtxt.APIClient.ServerVersion))

	// Entries
	passwordEntry := widget.NewPasswordEntry()
	usernameEntry := widget.NewEntry()
	usernameEntry.SetText(appCtxt.APIClient.CurrentUser.Username)
	emailEntry := widget.NewEntry()
	emailEntry.SetText(appCtxt.APIClient.CurrentUser.Email)
	// Group entries in a form
	passwordForm := widget.NewFormItem("Password", passwordEntry)
	usernameForm := widget.NewFormItem("Username", usernameEntry)
	emailForm := widget.NewFormItem("Email", emailEntry)
	contentForm := []*widget.FormItem{usernameForm, emailForm, passwordForm}

	// Special form for the dialog password confirmation window
	confirmPasswordform := []*widget.FormItem{passwordForm}

	// Buttons
	updateButton := widget.NewButtonWithIcon("Update info", theme.DocumentCreateIcon(), func() {
		buttonFuncUpdateUserInfo(appCtxt, confirmPasswordform, contentForm, passwordEntry, usernameEntry, emailEntry, statusLabel)
	})

	exitButton := widget.NewButtonWithIcon("Homepage", theme.HomeIcon(), func() {
		dialog.ShowConfirm("Exit", "Are you sure you want to go back to Homepage ?\nAll unsubmitted changes will be lost!", func(b bool) {
			if b {
				appCtxt.PageManager.ShowHomePage()
			}
		}, appCtxt.MainWindow)
	})

	deleteUserButton := widget.NewButtonWithIcon("Delete User", theme.DeleteIcon(), func() {
		buttonFuncDeleteUser(appCtxt)
	})

	// Group objects
	textColumn := container.NewVBox(layout.NewSpacer(), clientVersion, serverVersion, usernameLabel, emailLabel, layout.NewSpacer(), statusLabel, updateButton, customSpacerVertical(100), deleteUserButton, layout.NewSpacer())
	centerRow := container.NewHBox(layout.NewSpacer(), textColumn, layout.NewSpacer())

	// Create the global frame
	globalContainer := container.NewBorder(
		titleText,
		exitButton,
		customSpacerHorizontal(100),
		customSpacerHorizontal(100),
		centerRow,
	)
	return globalContainer

}

func buttonFuncUpdateUserInfo(appCtxt *context.AppContext, confirmPasswordform, contentForm []*widget.FormItem, passwordEntry, usernameEntry, emailEntry *widget.Entry, statusLabel *widget.Label) {
	dialog.ShowForm("Confirm your password", "Confirm", "Cancel", confirmPasswordform, func(b bool) {
		if b {
			if err := appCtxt.APIClient.Auth.ConfirmPassword(passwordEntry.Text); err != nil {
				statusLabel.SetText("Wrong password")
			} else {
				passwordEntry.SetText("")
				dialog.ShowForm("ALL info required, including not updated fields", "Confirm", "Cancel", contentForm, func(b bool) {
					if b {
						// Call the Update User client API function
						_, err := appCtxt.APIClient.Users.UpdateUser(usernameEntry.Text, passwordEntry.Text, emailEntry.Text)
						if err != nil {
							switch err {
							case models.ErrUnauthorized:
								if _, err2 := appCtxt.APIClient.Auth.RefreshTokens(); err2 != nil {
									dialog.ShowConfirm("Authorization problem", "There is a problem with your authorization,\nyou'll be redirected to Login page", func(b bool) {
										appCtxt.PageManager.ShowLoginPage()
									}, appCtxt.MainWindow)
								} else {
									dialog.ShowInformation("Information", "Client needed to refresh your acess token\nSorry for the inconvenience\nPlease try again, it should work now !", appCtxt.MainWindow)
								}
							case models.ErrServerIssue:
								dialog.ShowInformation("Error", "Error with server, please retry later", appCtxt.MainWindow)
							case models.ErrBadRequest:
								statusLabel.SetText("User's info update failed: one field was not provided")
							case models.ErrConflict:
								statusLabel.SetText("User's info update failed: username or email already used")
							default:
								dialog.ShowError(err, appCtxt.MainWindow)
							}
						} else {
							dialog.ShowConfirm("Information", "Info successfully updated !\nYou'll need to log in again", func(b bool) {
								if b {
									appCtxt.PageManager.ShowLoginPage()
								} else {
									appCtxt.MainWindow.Close()
								}
							}, appCtxt.MainWindow)
						}
					} else {
						passwordEntry.SetText("")
					}
				}, appCtxt.MainWindow)
			}
		}
	}, appCtxt.MainWindow)
}

func buttonFuncDeleteUser(appCtxt *context.AppContext) {
	dialog.ShowConfirm("Confirm", "Are you sure you want to delete your user account ??", func(b bool) {
		if b {
			dialog.ShowConfirm("Last warning", "Last warning, This action is irreversible !", func(b bool) {
				if b {
					err := appCtxt.APIClient.Users.DeleteUser()
					if err != nil {
						dialog.ShowError(err, appCtxt.MainWindow)
					}
					appCtxt.PageManager.ShowLoginPage()
				}
			}, appCtxt.MainWindow)
		}
	}, appCtxt.MainWindow)
}
