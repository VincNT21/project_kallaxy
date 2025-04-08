package gui

import (
	"errors"
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

	// ImageUrl row (with "get info online" button)
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
		if titleEntry.Text == "" {
			// If title not provided
			dialog.ShowInformation("Info", "You need to provide a title before clicking on this !", appCtxt.MainWindow)
		} else {
			// If title provided
			initSearchResultContent(appCtxt, appCtxt.MainWindow, titleEntry.Text, mediaType, "", func(selectedMedium models.ClientMedium) {
				updateFormWithSearchResult(selectedMedium, mediaForm)
				mediaForm.Refresh()
				updateMetadataForm(appCtxt, selectedMedium, metadataEntryMap)
				metadataForm.Refresh()
			})
		}

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
				// If Confirmed. call the CreateMediumAndRecord client API function
				if b {
					// PLACEHOLDER
					fmt.Println(metadataEntryMap)
					// PLACEHOLDER
					_, _, err := appCtxt.APIClient.Media.CreateMediumAndRecord(
						titleEntry.Text,
						mediaTypeEntry.Text,
						creatorEntry.Text,
						pubDateEntry.Text,
						imageUrlEntry.Text,
						startDateEntry.Text,
						endDateEntry.Text,
					)
					if err != nil {
						switch err {
						case models.ErrUnauthorized:
							if _, err2 := appCtxt.APIClient.Auth.RefreshTokens(); err2 != nil {
								dialog.NewConfirm("Authorization problem", "There is a problem with your authorization,\nyou'll be redirected to Login page", func(b bool) {
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
	})

	// metadataForm := widget.NewForm()
	groupForms := container.NewVBox(mediaForm, urlRow, widget.NewSeparator(), metadataForm, widget.NewSeparator(), recordForm)
	submitRow := container.NewBorder(nil, nil, customSpacerHorizontal(20), customSpacerHorizontal(20), submitButton)
	statusRow := container.NewHBox(layout.NewSpacer(), statusLabel, layout.NewSpacer())
	centralPart := container.NewVBox(groupForms, statusRow, submitRow)

	// Create the global frame
	globalContainer := container.NewBorder(pageTitleText, exitButton, nil, nil, centralPart)

	return globalContainer
}

func initSearchResultContent(appCtxt *context.AppContext, parentWindow fyne.Window, mediumTitle, mediumType, vgPlatform string, onConfirm func(models.ClientMedium)) {
	// Create the window
	secondaryWindow := fyne.CurrentApp().NewWindow("Search Results")
	secondaryWindow.CenterOnScreen()
	secondaryWindow.Resize(fyne.NewSize(640, 480))

	// Get results list
	results, err := appCtxt.APIClient.Helpers.SearchMediaOnExternalApiByTitle(mediumType, mediumTitle, vgPlatform)
	if err != nil {
		if err == models.ErrNotFound {
			dialog.ShowError(errors.New("no media found with this title"), parentWindow)
		} else {
			dialog.ShowError(fmt.Errorf("an error occured while trying to get search online results\n%v", err), parentWindow)
		}
		secondaryWindow.Close()
		return
	}

	// Initialize with the first result
	updateSearchResultContent(appCtxt, mediumType, secondaryWindow, results, 0, onConfirm)

	// Display window
	secondaryWindow.Show()
}

func updateSearchResultContent(appCtxt *context.AppContext, mediaType string, w fyne.Window, results []models.ShortOnlineSearchResult, i int, onConfirm func(models.ClientMedium)) {

	result := results[i]

	// Create UI components
	// Texts
	pageTitleText := canvas.NewText(fmt.Sprintf("Result %v / %v", result.Num, result.TotalNumFound), color.White)
	pageTitleText.TextSize = 20
	pageTitleText.Alignment = fyne.TextAlignCenter
	pageTitleText.TextStyle.Bold = true

	titleText := canvas.NewText(fmt.Sprintf("Title: %s", result.Title), color.White)
	titleText.TextSize = 16
	titleText.Alignment = fyne.TextAlignCenter
	titleText.TextStyle.Bold = true

	statusText := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})

	// Fetch the image as a buffer
	bufImage, err := appCtxt.APIClient.Helpers.GetImage(result.ImageUrl)
	if err != nil {
		dialog.ShowError(err, w)
		w.Close()
	}
	// Create the image component
	image := canvas.NewImageFromReader(bufImage, "image")
	image.FillMode = canvas.ImageFillContain
	image.Resize(fyne.NewSize(350, 250))

	// Buttons
	detailsButton := widget.NewButtonWithIcon("Get details", theme.SearchIcon(), func() {
		showSearchMediumDetails(appCtxt, mediaType, result.ApiID, result.ImageUrl, w, onConfirm)

	})
	cancelButton := widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
		w.Close()
	})
	nextButton := widget.NewButtonWithIcon("Next result", theme.NavigateNextIcon(), func() {
		if i+1 == result.TotalNumFound {
			statusText.SetText("This is the last result")
		} else {
			// Show next page of results
			updateSearchResultContent(appCtxt, mediaType, w, results, i+1, onConfirm)
		}
	})
	previousButton := widget.NewButtonWithIcon("Previous result", theme.NavigateBackIcon(), func() {
		if i == 0 {
			statusText.SetText("This is the first result")
		} else {
			// Show previous page of results
			updateSearchResultContent(appCtxt, mediaType, w, results, i-1, onConfirm)
		}
	})

	// Layout the elements
	globalContainer := container.NewBorder(
		pageTitleText, // Top
		container.NewVBox( // Bottom
			titleText,
			container.NewHBox(
				layout.NewSpacer(),
				previousButton,
				nextButton,
				layout.NewSpacer(),
			),
			statusText,
			container.NewHBox(
				cancelButton,
				layout.NewSpacer(),
				detailsButton,
			),
		),
		nil,   // Left
		nil,   // Right
		image, // Center
	)

	// Set container to window
	w.SetContent(globalContainer)
}

func showSearchMediumDetails(appCtxt *context.AppContext, mediaType, mediumApiID, imageUrl string, parentWindow fyne.Window, onConfirm func(models.ClientMedium)) {
	// Get details for medium on external API
	medium, err := appCtxt.APIClient.Helpers.SearchMediumDetailsOnExternalApi(mediaType, mediumApiID)
	if err != nil {
		dialog.ShowError(fmt.Errorf("couldn't get details about medium: %v", err), parentWindow)
		return
	}

	// Prepare results
	titleText := canvas.NewText(fmt.Sprintf("Title: %s", medium.Title), color.White)
	titleText.TextSize = 16
	creatorText := canvas.NewText(fmt.Sprintf("Creator: %s", medium.Creator), color.White)
	creatorText.TextSize = 16
	pubDateText := canvas.NewText(fmt.Sprintf("Publication date: %s", medium.PubDate), color.White)
	pubDateText.TextSize = 16
	// Add metadata ?

	// Display them in a dialog box
	dialog.ShowCustomConfirm(
		"Details",
		"Confirm",
		"Dismiss",
		container.NewVBox(titleText, creatorText, pubDateText),
		func(b bool) {
			if b {
				// If user confirms, call OnConfirm callback function
				medium.ImageUrl = imageUrl // Insert back the proper image url
				onConfirm(medium)
			}
		},
		parentWindow,
	)

}

func updateFormWithSearchResult(result models.ClientMedium, form *widget.Form) {
	for _, item := range form.Items {
		entry, isEntry := item.Widget.(*widget.Entry)
		if !isEntry {
			continue
		}
		// Update based on field name
		switch item.Text {
		case "Title":
			entry.SetText(result.Title)
		case "Creator":
			entry.SetText(result.Creator)
		case "Publication date":
			entry.SetText(result.PubDate)
		}

	}
}

func updateMetadataForm(appCtxt *context.AppContext, result models.ClientMedium, entryMap map[string]*widget.Entry) {
	// Iterate through each field/widget on the entry map
	for field, entryWidget := range entryMap {
		// Check if field exists in the metadata
		if metadataValue, exists := result.Metadata[field]; exists {
			// If yes, convert the interface{} value to string and set the entry
			stringValue := formatMetadataValueForEntry(appCtxt, field, metadataValue)
			entryWidget.SetText(stringValue)
		}
	}
}
