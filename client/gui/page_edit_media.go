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

func createMediaEditingContent(appCtxt *context.AppContext, mediaType string, mediumWithRecord models.MediumWithRecord) *fyne.Container {

	// Create UI objects
	// Texts
	statusLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})
	pageTitleText := canvas.NewText(fmt.Sprintf("Update a %s's info", mediaType), color.White)
	pageTitleText.TextSize = 20
	pageTitleText.Alignment = fyne.TextAlignCenter
	pageTitleText.TextStyle.Bold = true

	// Entries and forms
	// Global Media forms
	titleEntry := widget.NewEntry()
	titleForm := widget.NewFormItem("Title", titleEntry)

	mediaTypeEntry := widget.NewEntry()
	mediaTypeEntry.SetText(mediaType)
	mediaTypeEntry.Disable()
	mediaTypeForm := widget.NewFormItem("Media Type", mediaTypeEntry)

	creatorEntry := widget.NewEntry()
	creatorForm := widget.NewFormItem("Creator", creatorEntry)

	pubDateEntry := widget.NewEntry()
	pubDateForm := widget.NewFormItem("Publication date", pubDateEntry)

	mediaForm := widget.NewForm(titleForm, mediaTypeForm, creatorForm, pubDateForm)

	imageUrlEntry := widget.NewEntry()
	imageUrlForm := widget.NewForm(widget.NewFormItem("Image URL", imageUrlEntry))

	// Metadata formItems will depend on the media_type
	metadataForm, metadataEntryMap := createMetadataForm(appCtxt, mediaType)

	// Set initial text in entries according to result
	setInitialValues(appCtxt, titleEntry, creatorEntry, pubDateEntry, imageUrlEntry, mediumWithRecord, metadataEntryMap)

	// UI Buttons
	// Undo changes Button (to get back to existing data)
	buttonUndoChanges := widget.NewButtonWithIcon("Undo all changes", theme.ContentUndoIcon(), func() {
		dialog.ShowConfirm("Confirm", "Are you sure you want to undo all changes ?", func(b bool) {
			if b {
				setInitialValues(appCtxt, titleEntry, creatorEntry, pubDateEntry, imageUrlEntry, mediumWithRecord, metadataEntryMap)
			}
		}, appCtxt.MainWindow)
	})

	mediumIdFromDB := widget.NewEntry()

	// Get info online button (positionned right to ImagURL form)
	buttonGetInfoOnline := widget.NewButtonWithIcon("Get Info Online\nfrom title", theme.DownloadIcon(), func() {
		buttonFuncGetInfoOnline(appCtxt, mediaType, mediaForm, imageUrlForm, metadataForm, imageUrlEntry, titleEntry, mediumIdFromDB, metadataEntryMap)
	})
	urlRow := container.NewBorder(nil, nil, nil, buttonGetInfoOnline, imageUrlForm)

	exitButton := widget.NewButtonWithIcon("Homepage", theme.HomeIcon(), func() {
		dialog.ShowConfirm("Exit", "Are you sure you want to go back to Homepage ?\n\nAll unsubmitted changes will be lost!", func(b bool) {
			if b {
				appCtxt.PageManager.ShowHomePage()
			}
		}, appCtxt.MainWindow)
	})

	submitButton := widget.NewButtonWithIcon("Update", theme.ConfirmIcon(), func() {
		buttonFuncSubmitUpdate(appCtxt, mediumWithRecord, mediaTypeEntry, titleEntry, creatorEntry, pubDateEntry, imageUrlEntry, metadataEntryMap)
	})

	// Group objects
	groupForms := container.NewVBox(mediaForm, urlRow, widget.NewSeparator(), metadataForm)
	submitRow := container.NewBorder(nil, nil, customSpacerHorizontal(20), customSpacerHorizontal(20), submitButton)
	statusRow := container.NewHBox(layout.NewSpacer(), statusLabel, layout.NewSpacer())
	centralPart := container.NewVBox(groupForms, statusRow, submitRow)

	// Create the global frame
	globalContainer := container.NewBorder(
		container.NewBorder(nil, nil, nil, buttonUndoChanges, pageTitleText),
		exitButton,
		nil, nil,
		centralPart,
	)

	return globalContainer
}

func setInitialValues(appCtxt *context.AppContext, titleEntry, creatorEntry, pubDateEntry, imageUrlEntry *widget.Entry, mediumWithRecord models.MediumWithRecord, metadataEntryMap map[string]*widget.Entry) {
	titleEntry.SetText(mediumWithRecord.Title)
	creatorEntry.SetText(mediumWithRecord.Creator)
	pubDateEntry.SetText(mediumWithRecord.PubDate)
	imageUrlEntry.SetText(mediumWithRecord.ImageUrl)
	updateMetadataForm(appCtxt, mediumWithRecord.Metadata, metadataEntryMap)
}

func buttonFuncSubmitUpdate(appCtxt *context.AppContext, mediumWithRecord models.MediumWithRecord, mediaTypeEntry, titleEntry, creatorEntry, pubDateEntry, imageUrlEntry *widget.Entry, metadataEntryMap map[string]*widget.Entry) {
	// Confirm info dialog box
	dialog.ShowCustomConfirm(
		"Confirm",
		"Update",
		"Cancel",
		container.NewVBox(
			widget.NewLabelWithStyle("Do you confirm the info entered ?", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle(fmt.Sprintf("Media Type: %s", mediaTypeEntry.Text), fyne.TextAlignLeading, fyne.TextStyle{}),
			widget.NewLabelWithStyle(fmt.Sprintf("Title: %s", titleEntry.Text), fyne.TextAlignLeading, fyne.TextStyle{}),
			widget.NewLabelWithStyle(fmt.Sprintf("Creator: %s", creatorEntry.Text), fyne.TextAlignLeading, fyne.TextStyle{}),
			widget.NewLabelWithStyle(fmt.Sprintf("Publication date: %s", pubDateEntry.Text), fyne.TextAlignLeading, fyne.TextStyle{}),
			widget.NewLabelWithStyle(fmt.Sprintf("Image URL: %s", imageUrlEntry.Text), fyne.TextAlignLeading, fyne.TextStyle{}),
		),
		func(b bool) {
			// If Confirmed. call the UpdateMedium client API function
			if b {
				// Get proper metadata field (according to fields specs)
				metadataParsed := extractMetadataValues(appCtxt, metadataEntryMap)

				// Make request to server
				_, err := appCtxt.APIClient.Media.UpdateMedium(
					mediumWithRecord.MediaID,
					titleEntry.Text,
					creatorEntry.Text,
					pubDateEntry.Text,
					imageUrlEntry.Text,
					metadataParsed,
				)
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
						dialog.ShowInformation("Error", "There is a problem with your request:\n- One field is missing in the form\nAND/OR\n- Start date is before end date\nPlease verify all fields", appCtxt.MainWindow)
					case models.ErrConflict:
						dialog.ShowInformation("Error", "A medium with the same couple title & media type already exists", appCtxt.MainWindow)
					case models.ErrNotFound:
						dialog.ShowInformation("Error", "Not found", appCtxt.MainWindow)
					default:
						dialog.ShowError(err, appCtxt.MainWindow)
					}
				} else {
					log.Println("--GUI-- UpdateMedium Form successful")
					dialog.ShowInformation("Updated", "Media update successful !", appCtxt.MainWindow)
					appCtxt.PageManager.ShowHomePage()
				}
			}
		}, appCtxt.MainWindow,
	)
}
