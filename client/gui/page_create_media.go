package gui

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/VincNT21/kallaxy/client/context"
	"github.com/VincNT21/kallaxy/client/models"
	datepicker "github.com/sdassow/fyne-datepicker"
)

func createMediaCreationContent(appCtxt *context.AppContext, mediaType string) *fyne.Container {
	// This variable will hold the mediumID, in case of a successful search on server's database
	mediumIdFromDB := widget.NewEntry()

	// Create UI objects
	// Texts
	statusLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})
	pageTitleText := canvas.NewText(fmt.Sprintf("Create a new %s", mediaType), color.White)
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

	// User's Record forms
	// Date entries have a Action Button that calls a Date Picker dialog
	startDateEntry := widget.NewEntry()
	startDateEntry.SetPlaceHolder("2025/01/01")
	startDateEntry.ActionItem = widget.NewButtonWithIcon("", theme.MoreHorizontalIcon(), func() {
		when := time.Now().UTC()

		if startDateEntry.Text != "" {
			t, err := time.Parse("2006/01/02", startDateEntry.Text)
			if err == nil {
				when = t
			}
		}

		datepicker := datepicker.NewDatePicker(when, time.Monday, func(when time.Time, ok bool) {
			if ok {
				startDateEntry.SetText(when.Format("2006/01/02"))
			}
		})

		dialog.ShowCustomConfirm(
			"Choose date",
			"Ok",
			"Cancel",
			datepicker,
			datepicker.OnActioned,
			appCtxt.MainWindow,
		)
	})
	startDateFormItem := widget.NewFormItem("Started on", startDateEntry)

	endDateEntry := widget.NewEntry()
	endDateEntry.SetPlaceHolder("2025/01/01")
	endDateEntry.ActionItem = widget.NewButtonWithIcon("", theme.MoreHorizontalIcon(), func() {
		when := time.Now().UTC()

		if endDateEntry.Text != "" {
			t, err := time.Parse("2006/01/02", endDateEntry.Text)
			if err == nil {
				when = t
			}
		}

		datepicker := datepicker.NewDatePicker(when, time.Monday, func(when time.Time, ok bool) {
			if ok {
				endDateEntry.SetText(when.Format("2006/01/02"))
			}
		})

		dialog.ShowCustomConfirm(
			"Choose date",
			"Ok",
			"Cancel",
			datepicker,
			datepicker.OnActioned,
			appCtxt.MainWindow,
		)
	})
	endDateFormItem := widget.NewFormItem("Completed on", endDateEntry)

	commentsEntry := widget.NewMultiLineEntry()
	commentsFormItem := widget.NewFormItem("Personal comments", commentsEntry)

	recordForm := widget.NewForm(startDateFormItem, endDateFormItem, commentsFormItem)

	// Metadata formItems will depend on the media_type
	metadataForm, metadataEntryMap := createMetadataForm(appCtxt, mediaType)

	// UI Buttons
	// Get info online button ("linked" to ImagURL form)
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

	submitButton := widget.NewButtonWithIcon("Submit", theme.ConfirmIcon(), func() {
		buttonFuncSubmitCreateMedia(appCtxt, mediaTypeEntry, titleEntry, creatorEntry, pubDateEntry, imageUrlEntry, startDateEntry, endDateEntry, mediumIdFromDB, metadataEntryMap)
	})

	// Group objects
	groupForms := container.NewVBox(mediaForm, urlRow, widget.NewSeparator(), metadataForm, widget.NewSeparator(), recordForm)
	submitRow := container.NewBorder(nil, nil, customSpacerHorizontal(20), customSpacerHorizontal(20), submitButton)
	statusRow := container.NewHBox(layout.NewSpacer(), statusLabel, layout.NewSpacer())
	centralPart := container.NewVBox(groupForms, statusRow, submitRow)

	// Create the global frame
	globalContainer := container.NewBorder(pageTitleText, exitButton, nil, nil, centralPart)

	return globalContainer
}

func buttonFuncGetInfoOnline(appCtxt *context.AppContext, mediaType string, mediaForm, imageUrlForm, metadataForm *widget.Form, imageUrlEntry, titleEntry, mediumIdFromDB *widget.Entry, metadataEntryMap map[string]*widget.Entry) {
	// First, check if couple mediaTitle+mediaType already exists in server's DB
	dbMedium, err := appCtxt.APIClient.Helpers.SearchMediumInDB(mediaType, titleEntry.Text)
	if err == nil {
		// Means we found medium in server's db
		// Inform user
		dialog.ShowInformation("Info", "This medium has been found on the server's database !\nYou can add your personal record info\nBut any update to medium's info will have no effects\nYou can do that later by going to your Shelf", appCtxt.MainWindow)
		// Update data
		updateFormWithSearchResult(dbMedium.Title, dbMedium.Creator, dbMedium.PubDate, mediaForm)
		mediaForm.Refresh()
		imageUrlEntry.SetText(dbMedium.ImageUrl)
		imageUrlForm.Refresh()
		updateMetadataForm(appCtxt, dbMedium.Metadata, metadataEntryMap)
		metadataForm.Refresh()
		// Disable fields modification (except for records)
		mediaForm.Disable()
		imageUrlForm.Disable()
		metadataForm.Disable()
		// Set the mediumIdFromDB variable with dbMedium's ID
		mediumIdFromDB.SetText(dbMedium.ID)
		return
	}

	// Otherwise, continue with online search

	// Book is a special media type, it's better to search by ISBN rather than by title
	if mediaType == "book" {
		// In case of book, ask for ISBN
		titleLine1 := canvas.NewText("Book search is more efficient with provided ISBN", color.White)
		titleLine1.Alignment = fyne.TextAlignCenter
		titleLine1.TextSize = 14
		titleLine2 := canvas.NewText("If you have one, please enter it", color.White)
		titleLine2.Alignment = fyne.TextAlignCenter
		titleLine2.TextSize = 14
		titleLine3 := canvas.NewText("or continue with search by Title", color.White)
		titleLine3.Alignment = fyne.TextAlignCenter
		titleLine3.TextSize = 14

		isbnEntry := widget.NewEntry()
		isbnEntry.SetPlaceHolder("ISBN 10 or ISBN 13, numbers only")

		dialog.ShowCustomConfirm(
			"Book search",
			"Search with ISBN",
			"Search with title only",
			container.NewVBox(titleLine1, titleLine2, titleLine3, layout.NewSpacer(), isbnEntry, layout.NewSpacer()),
			func(b bool) {
				if b {
					showSearchMediumDetails(appCtxt, "book", isbnEntry.Text, "", appCtxt.MainWindow, metadataEntryMap, func(selectedMedium models.ClientMedium) {
						// OnConfirm function
						updateFormWithSearchResult(selectedMedium.Title, selectedMedium.Creator, selectedMedium.PubDate, mediaForm)
						mediaForm.Refresh()
						imageUrlEntry.SetText(selectedMedium.ImageUrl)
						imageUrlForm.Refresh()
						updateMetadataForm(appCtxt, selectedMedium.Metadata, metadataEntryMap)
						metadataForm.Refresh()
					})
				} else {
					if titleEntry.Text == "" {
						// If title not provided
						dialog.ShowInformation("Info", "You need to provide a title before clicking on this !", appCtxt.MainWindow)
						return
					}
					initSearchResultContent(appCtxt, appCtxt.MainWindow, titleEntry.Text, mediaType, "", metadataEntryMap, func(selectedMedium models.ClientMedium) {
						// OnConfirm function
						updateFormWithSearchResult(selectedMedium.Title, selectedMedium.Creator, selectedMedium.PubDate, mediaForm)
						mediaForm.Refresh()
						imageUrlEntry.SetText(selectedMedium.ImageUrl)
						imageUrlForm.Refresh()
						updateMetadataForm(appCtxt, selectedMedium.Metadata, metadataEntryMap)
						metadataForm.Refresh()
					})
				}
			},
			appCtxt.MainWindow,
		)
		return
	}

	// Normal case (not book)
	if titleEntry.Text == "" {
		// If title not provided
		dialog.ShowInformation("Info", "You need to provide a title before clicking on this !", appCtxt.MainWindow)
		return
	} else {
		// Search function
		initSearchResultContent(appCtxt, appCtxt.MainWindow, titleEntry.Text, mediaType, "", metadataEntryMap, func(selectedMedium models.ClientMedium) {
			// OnConfirm function
			updateFormWithSearchResult(selectedMedium.Title, selectedMedium.Creator, selectedMedium.PubDate, mediaForm)
			mediaForm.Refresh()
			imageUrlEntry.SetText(selectedMedium.ImageUrl)
			imageUrlForm.Refresh()
			updateMetadataForm(appCtxt, selectedMedium.Metadata, metadataEntryMap)
			metadataForm.Refresh()
		})
	}
}

func buttonFuncSubmitCreateMedia(appCtxt *context.AppContext, mediaTypeEntry, titleEntry, creatorEntry, pubDateEntry, imageUrlEntry, startDateEntry, endDateEntry, mediumIdFromDB *widget.Entry, metadataEntryMap map[string]*widget.Entry) {
	// Confirm info dialog box
	dialog.ShowCustomConfirm(
		"Confirm",
		"Create",
		"Cancel",
		container.NewVBox(
			widget.NewLabelWithStyle("Do you confirm the info entered ?", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle(fmt.Sprintf("Media Type: %s", mediaTypeEntry.Text), fyne.TextAlignLeading, fyne.TextStyle{}),
			widget.NewLabelWithStyle(fmt.Sprintf("Title: %s", titleEntry.Text), fyne.TextAlignLeading, fyne.TextStyle{}),
			widget.NewLabelWithStyle(fmt.Sprintf("Creator: %s", creatorEntry.Text), fyne.TextAlignLeading, fyne.TextStyle{}),
			widget.NewLabelWithStyle(fmt.Sprintf("Publication date: %s", pubDateEntry.Text), fyne.TextAlignLeading, fyne.TextStyle{}),
			widget.NewLabelWithStyle(fmt.Sprintf("Image URL: %s", imageUrlEntry.Text), fyne.TextAlignLeading, fyne.TextStyle{}),
			widget.NewLabelWithStyle(fmt.Sprintf("Start Date: %s", startDateEntry.Text), fyne.TextAlignLeading, fyne.TextStyle{}),
			widget.NewLabelWithStyle(fmt.Sprintf("End Date: %s", endDateEntry.Text), fyne.TextAlignLeading, fyne.TextStyle{}),
		),
		func(b bool) {

			if b {
				// First, check if medium's info comes from DB or not, by looking at mediumIdFromDB Entry value
				if mediumIdFromDB.Text != "" {
					_, err := appCtxt.APIClient.Records.CreateRecord(mediumIdFromDB.Text, startDateEntry.Text, endDateEntry.Text)
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
							dialog.ShowInformation("Error", "There is a problem with your request:\n Start date is before end date", appCtxt.MainWindow)
						case models.ErrConflict:
							dialog.ShowInformation("Error", "You already have a personal record for this media", appCtxt.MainWindow)
						case models.ErrNotFound:
							dialog.ShowInformation("Error", "Not found", appCtxt.MainWindow)
						default:
							dialog.ShowError(err, appCtxt.MainWindow)
						}
					} else {
						log.Println("--GUI-- CreateMedia Form successful")
						dialog.ShowInformation("Created", "Personal record creation successful !", appCtxt.MainWindow)
						appCtxt.PageManager.ShowHomePage()
					}
					return
				}

				// If not, call the CreateMediumAndRecord client API function
				// Get proper metadata field (according to fields specs)
				metadataParsed := extractMetadataValues(appCtxt, metadataEntryMap)

				// Make request to server
				_, _, err := appCtxt.APIClient.Media.CreateMediumAndRecord(
					titleEntry.Text,
					mediaTypeEntry.Text,
					creatorEntry.Text,
					pubDateEntry.Text,
					imageUrlEntry.Text,
					startDateEntry.Text,
					endDateEntry.Text,
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
					log.Println("--GUI-- CreateMedia Form successful")
					dialog.ShowInformation("Created", "Media creation successful !", appCtxt.MainWindow)
					appCtxt.PageManager.ShowHomePage()
				}
			}
		}, appCtxt.MainWindow,
	)
}

func updateFormWithSearchResult(resultTitle, resultCreator, resultPubdate string, form *widget.Form) {
	for _, item := range form.Items {
		entry, isEntry := item.Widget.(*widget.Entry)
		if !isEntry {
			continue
		}
		// Update based on field name
		switch item.Text {
		case "Title":
			entry.SetText(resultTitle)
		case "Creator":
			entry.SetText(resultCreator)
		case "Publication date":
			entry.SetText(resultPubdate)
		}

	}
}

func updateMetadataForm(appCtxt *context.AppContext, resultMetadata map[string]interface{}, entryMap map[string]*widget.Entry) {
	// Iterate through each field/widget on the entry map
	for field, entryWidget := range entryMap {
		// Check if field exists in the metadata
		if metadataValue, exists := resultMetadata[field]; exists {
			// If yes, convert the interface{} value to string and set the entry
			stringValue := formatMetadataValueForEntry(appCtxt, field, metadataValue)
			entryWidget.SetText(stringValue)
		}
	}
}
