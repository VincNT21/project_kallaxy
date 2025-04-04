package gui

import (
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

func (pm *GuiPageManager) GetHomeWindow() {
	// Create the window
	w := pm.appGui.NewWindow("Kallaxy")
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(800, 600))

	// Create objects
	titleText := canvas.NewText("Welcome to your Kallaxy", color.White)
	titleText.TextSize = 18
	titleText.Alignment = fyne.TextAlignCenter
	titleText.TextStyle.Bold = true
	usernameText := canvas.NewText(pm.appCtxt.APIClient.CurrentUser.Username, color.White)
	usernameText.TextSize = 20
	usernameText.Alignment = fyne.TextAlignCenter
	usernameText.TextStyle.Bold = true

	addMediaButton := widget.NewButton("Add New Media", func() {
		// Create the custom dialog box to choose media_type
		mediaTypeQuestion := canvas.NewText("Which type of media ?", color.White)

		// Create a button for each media Type + an "other" button
		buttonBook := widget.NewButtonWithIcon("Book", theme.DocumentIcon(), func() {
			pm.mediaType = "book"
			pm.GetCreateMediaWindow()
			w.Close()
		})
		buttonMovie := widget.NewButtonWithIcon("Movie", theme.MediaVideoIcon(), func() {
			pm.mediaType = "movie"
			pm.GetCreateMediaWindow()
			w.Close()
		})
		buttonSeries := widget.NewButtonWithIcon("Series", theme.MediaPlayIcon(), func() {
			pm.mediaType = "series"
			pm.GetCreateMediaWindow()
			w.Close()
		})
		buttonVideogame := widget.NewButtonWithIcon("Videogame", theme.ComputerIcon(), func() {
			pm.mediaType = "videogame"
			pm.GetCreateMediaWindow()
			w.Close()
		})
		buttonBoardgame := widget.NewButtonWithIcon("Boardgame", theme.SettingsIcon(), func() {
			pm.mediaType = "boardgame"
			pm.GetCreateMediaWindow()
			w.Close()
		})
		buttonOther := widget.NewButtonWithIcon("Other", theme.MoreHorizontalIcon(), func() {
			otherInput := widget.NewEntry()
			otherForm := widget.NewFormItem("Media Type", otherInput)
			dialog.ShowForm("Other media type", "Confirm", "Dismiss", []*widget.FormItem{otherForm}, func(b bool) {
				if b {
					pm.mediaType = otherInput.Text
					pm.GetCreateMediaWindow()
					w.Close()
				}
			}, w)
		})

		// Groups buttons
		groupButtons := container.NewVBox(buttonBook, buttonMovie, buttonSeries, buttonVideogame, buttonBoardgame, buttonOther)
		globalContainer := container.NewHBox(layout.NewSpacer(), groupButtons, layout.NewSpacer())

		// Display custom dialog
		dialog.ShowCustomWithoutButtons("Kallaxy", container.NewBorder(mediaTypeQuestion, nil, nil, nil, globalContainer), w)

	})
	showShelfButton := widget.NewButton("Show My Shelf", func() {
		pm.GetShelfWindow()
	})
	manageButton := widget.NewButtonWithIcon("Manage\nUser Parameters", theme.AccountIcon(), func() {
		pm.GetUserParametersWindow()
		w.Close()
	})

	logoutButton := widget.NewButtonWithIcon("Logout", theme.LogoutIcon(), func() {
		err := pm.appCtxt.APIClient.Auth.LogoutUser()
		if err != nil {
			log.Println("--GUI-- User failed to logout")
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
				dialog.ShowInformation("Logout error", "Error with server, please retry later", w)
			case models.ErrNotFound:
				dialog.ShowError(err, w)
			default:
				dialog.ShowError(err, w)
			}
		} else {
			dialog.ShowInformation("Logout successful", "You've been logout !", w)
			pm.GetLoginWindow()
			w.Close()
		}
	})
	exitButton := widget.NewButtonWithIcon("Exit App", theme.CancelIcon(), func() {
		dialog.ShowConfirm("Exit", "Are you sure you want to exit Kallaxy App ?", func(b bool) {
			if b {
				w.Close()
			}
		}, w)
	})

	// Create rows
	centralbuttonsRow := container.NewVBox(addMediaButton, &layout.Spacer{FixVertical: true}, showShelfButton)
	centralRow := container.NewHBox(layout.NewSpacer(), centralbuttonsRow, layout.NewSpacer())
	exitButtons := container.NewVBox(logoutButton, exitButton)
	bottomRow := container.NewHBox(manageButton, layout.NewSpacer(), exitButtons)
	globalContainer := container.NewVBox(layout.NewSpacer(), titleText, usernameText, layout.NewSpacer(), centralRow, layout.NewSpacer(), bottomRow)

	// Set container to window
	w.SetContent(globalContainer)
	w.Show()
}
