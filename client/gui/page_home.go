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
	"github.com/VincNT21/kallaxy/client/context"
	"github.com/VincNT21/kallaxy/client/models"
)

func createHomepageContent(appCtxt *context.AppContext) *fyne.Container {
	// Create objects
	// Top Texts
	titleText := canvas.NewText("Welcome to your Kallaxy", color.White)
	titleText.TextSize = 40
	titleText.Alignment = fyne.TextAlignCenter
	titleText.TextStyle.Bold = true

	usernameText := canvas.NewText(appCtxt.APIClient.CurrentUser.Username, color.White)
	usernameText.TextSize = 40
	usernameText.Alignment = fyne.TextAlignCenter
	usernameText.TextStyle.Bold = true

	// Buttons
	addMediaButton := widget.NewButton("Add New Media", func() {
		buttonFuncMediaTypeChoice(appCtxt)
	})

	showShelfButton := widget.NewButton("Show My Shelf", func() {
		appCtxt.PageManager.ShowShelfPage()
	})

	manageButton := widget.NewButtonWithIcon("Manage\nUser Parameters", theme.AccountIcon(), func() {
		appCtxt.PageManager.ShowParametersPage()
	})

	logoutButton := widget.NewButtonWithIcon("Logout", theme.LogoutIcon(), func() {
		buttonFuncLogout(appCtxt)
	})

	exitButton := widget.NewButtonWithIcon("Exit App", theme.CancelIcon(), func() {
		dialog.ShowConfirm("Exit", "Are you sure you want to exit Kallaxy App ?", func(b bool) {
			if b {
				appCtxt.MainWindow.Close()
			}
		}, appCtxt.MainWindow)
	})

	// Create rows
	centralbuttonsRow := container.NewVBox(addMediaButton, &layout.Spacer{FixVertical: true}, showShelfButton)
	centralRow := container.NewBorder(
		nil,
		nil,
		customSpacerHorizontal(100),
		customSpacerHorizontal(100),
		centralbuttonsRow,
	)
	exitButtons := container.NewVBox(logoutButton, exitButton)
	bottomRow := container.NewHBox(manageButton, layout.NewSpacer(), exitButtons)

	// Set the global frame container
	globalContainer := container.NewVBox(layout.NewSpacer(), titleText, usernameText, layout.NewSpacer(), centralRow, layout.NewSpacer(), bottomRow)

	// Return global container
	return globalContainer
}

func buttonFuncMediaTypeChoice(appCtxt *context.AppContext) {
	// Create a custom dialog box to choose media_type
	var mediaTypeDialog dialog.Dialog

	mediaTypeQuestion := canvas.NewText("Which type of media ?", color.White)
	mediaTypeQuestion.Alignment = fyne.TextAlignCenter
	mediaTypeQuestion.TextSize = 20

	// Create a button for each media Type + an "other" button
	buttonBook := widget.NewButtonWithIcon("Book", theme.DocumentIcon(), func() {
		appCtxt.PageManager.ShowCreateMediaPage("book")
		mediaTypeDialog.Hide()
	})
	buttonMovie := widget.NewButtonWithIcon("Movie", theme.MediaVideoIcon(), func() {
		appCtxt.PageManager.ShowCreateMediaPage("movie")
		mediaTypeDialog.Hide()
	})
	buttonSeries := widget.NewButtonWithIcon("Series", theme.MediaPlayIcon(), func() {
		appCtxt.PageManager.ShowCreateMediaPage("series")
		mediaTypeDialog.Hide()
	})
	buttonVideogame := widget.NewButtonWithIcon("Videogame", theme.ComputerIcon(), func() {
		appCtxt.PageManager.ShowCreateMediaPage("videogame")
		mediaTypeDialog.Hide()
	})
	buttonBoardgame := widget.NewButtonWithIcon("Boardgame", theme.SettingsIcon(), func() {
		appCtxt.PageManager.ShowCreateMediaPage("boardgame")
		mediaTypeDialog.Hide()
	})
	buttonOther := widget.NewButtonWithIcon("Other", theme.MoreHorizontalIcon(), func() {
		otherInput := widget.NewEntry()
		otherForm := widget.NewFormItem("Media Type", otherInput)
		dialog.ShowForm("Other media type", "Confirm", "Dismiss", []*widget.FormItem{otherForm}, func(b bool) {
			if b {
				appCtxt.PageManager.ShowCreateMediaPage(otherInput.Text)
				mediaTypeDialog.Hide()
			}
		}, appCtxt.MainWindow)
	})

	// Groups buttons
	groupButtons := container.NewVBox(buttonBook, buttonMovie, buttonSeries, buttonVideogame, buttonBoardgame, buttonOther)
	globalContainer := container.NewBorder(
		nil,
		nil,
		customSpacerHorizontal(50),
		customSpacerHorizontal(50),
		groupButtons,
	)

	// Set content and display custom dialog
	mediaTypeDialog = dialog.NewCustomWithoutButtons("", container.NewBorder(
		container.NewVBox(mediaTypeQuestion),
		nil,
		nil,
		nil,
		globalContainer,
	), appCtxt.MainWindow)

	mediaTypeDialog.Show()
}

func buttonFuncLogout(appCtxt *context.AppContext) {
	err := appCtxt.APIClient.Auth.LogoutUser()
	if err != nil {
		log.Println("--GUI-- User failed to logout")
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
			dialog.ShowInformation("Logout error", "Error with server, please retry later", appCtxt.MainWindow)
		case models.ErrNotFound:
			dialog.ShowError(err, appCtxt.MainWindow)
		default:
			dialog.ShowError(err, appCtxt.MainWindow)
		}
	} else {
		dialog.ShowInformation("Logout successful", "You've been logout !", appCtxt.MainWindow)
		appCtxt.PageManager.ShowLoginPage()
	}
}
