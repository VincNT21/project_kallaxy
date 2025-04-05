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
	"github.com/VincNT21/kallaxy/client/models"
)

func (pm *GuiPageManager) GetUserParametersWindow() {
	// Create the window
	w := pm.appGui.NewWindow("User Informations")
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(400, 300))

	// Create UI objects
	// Texts
	statusLabel := widget.NewLabelWithStyle("You'll be automatically logged out if you update your personal info", fyne.TextAlignCenter, fyne.TextStyle{})
	titleText := canvas.NewText("Personal Informations", color.White)
	titleText.TextSize = 20
	titleText.Alignment = fyne.TextAlignCenter
	titleText.TextStyle.Bold = true

	usernameLabel := widget.NewLabel(fmt.Sprintf("Username: %s", pm.appCtxt.APIClient.CurrentUser.Username))
	emailLabel := widget.NewLabel(fmt.Sprintf("Email: %s", pm.appCtxt.APIClient.CurrentUser.Email))

	// Entries
	passswordEntry := widget.NewPasswordEntry()
	usernameEntry := widget.NewEntry()
	usernameEntry.SetText(pm.appCtxt.APIClient.CurrentUser.Username)
	emailEntry := widget.NewEntry()
	emailEntry.SetText(pm.appCtxt.APIClient.CurrentUser.Email)
	// Group entries in a form
	passwordForm := widget.NewFormItem("Password", passswordEntry)
	usernameForm := widget.NewFormItem("Username", usernameEntry)
	emailForm := widget.NewFormItem("Email", emailEntry)
	contentForm := []*widget.FormItem{usernameForm, emailForm, passwordForm}

	// Special form for the dialog password confirmation window
	confirmPasswordform := []*widget.FormItem{passwordForm}

	// Buttons
	updateButton := widget.NewButtonWithIcon("Update info", theme.DocumentCreateIcon(), func() {
		dialog.ShowForm("Confirm your password", "Confirm", "Cancel", confirmPasswordform, func(b bool) {
			if err := pm.appCtxt.APIClient.Auth.ConfirmPassword(passswordEntry.Text); err != nil {
				statusLabel.SetText("Wrong password")
			} else {
				passswordEntry.SetText("")
				dialog.ShowForm("ALL info required, including not updated fields", "Confirm", "Cancel", contentForm, func(b bool) {
					if b {
						// Call the Update User client API function
						_, err := pm.appCtxt.APIClient.Users.UpdateUser(usernameEntry.Text, passswordEntry.Text, emailEntry.Text)
						if err != nil {
							switch err {
							case models.ErrUnauthorized:
								if _, err2 := pm.appCtxt.APIClient.Auth.RefreshTokens(); err2 != nil {
									dialog.NewConfirm("Authorization problem", "There is a problem with your authorization,\nyou'll be redirected to Login page", func(b bool) {
										pm.GetLoginWindow()
										w.Close()
									}, w)
								} else {
									dialog.ShowInformation("Information", "Client needed to refresh your acess token\nSorry for the inconvenience\nPlease try again, it should work now !", w)
								}
							case models.ErrServerIssue:
								dialog.ShowInformation("Error", "Error with server, please retry later", w)
							case models.ErrBadRequest:
								statusLabel.SetText("User's info update failed: one field was not provided")
							case models.ErrConflict:
								statusLabel.SetText("User's info update failed: username or email already used")
							default:
								dialog.ShowError(err, w)
							}
						} else {
							dialog.ShowConfirm("Information", "Info successfully updated !\nYou'll need to log in again", func(b bool) {
								if b {
									pm.GetLoginWindow()
									w.Close()
								} else {
									w.Close()
								}
							}, w)
						}
					} else {
						passswordEntry.SetText("")
					}
				}, w)
			}
		}, w)
	})

	exitButton := widget.NewButtonWithIcon("Homepage", theme.HomeIcon(), func() {
		dialog.ShowConfirm("Exit", "Are you sure you want to go back to Homepage ?\nAll unsubmitted changes will be lost!", func(b bool) {
			if b {
				pm.GetHomeWindow()
				w.Close()
			}
		}, w)
	})

	// Group objects
	textColumn := container.NewVBox(usernameLabel, layout.NewSpacer(), emailLabel)
	centerRow := container.NewHBox(layout.NewSpacer(), textColumn, layout.NewSpacer())
	bottomRow := container.NewHBox(updateButton, layout.NewSpacer(), exitButton)

	// Create the global frame
	globalContainer := container.NewVBox(layout.NewSpacer(), titleText, layout.NewSpacer(), centerRow, layout.NewSpacer(), statusLabel, layout.NewSpacer(), bottomRow)

	// Set container to window
	w.SetContent(globalContainer)
	w.Show()
}
