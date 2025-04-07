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

	// Get info online button
	buttonGetInfoOnline := widget.NewButtonWithIcon("Get Info Online\nfrom title", theme.DownloadIcon(), func() {
		if titleEntry.Text == "" {
			// If title not provided
			dialog.ShowInformation("Info", "You need to provide a title before clicking on this !", appCtxt.MainWindow)
		} else {
			// If title provided
			initSearchResultContent(appCtxt, appCtxt.MainWindow, titleEntry.Text, mediaType, "", func(s string) {})
		}

	})
	urlRow := container.NewBorder(nil, nil, nil, buttonGetInfoOnline, imageUrlForm)

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
	startDateForm := widget.NewFormItem("Started on", startDateEntry)

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
	endDateForm := widget.NewFormItem("Completed on", endDateEntry)

	commentsEntry := widget.NewMultiLineEntry()
	commentsForm := widget.NewFormItem("Personal comments", commentsEntry)

	recordForm := widget.NewForm(startDateForm, endDateForm, commentsForm)

	// Metadata formItems will depend on the media_type

	metadataForm, metadataEntryMap := createMetadataForm(appCtxt, mediaType)

	// UI Buttons
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

func initSearchResultContent(appCtxt *context.AppContext, parentWindow fyne.Window, mediumTitle, mediumType, vgPlatform string, onConfirm func(string)) {
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
	updateSearchResultContent(appCtxt, secondaryWindow, results, 0, onConfirm)

	// Display window
	secondaryWindow.Show()
}

func updateSearchResultContent(appCtxt *context.AppContext, w fyne.Window, results []models.ShortOnlineSearchResult, i int, onConfirm func(string)) {

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
	confirmButton := widget.NewButtonWithIcon("Confirm", theme.ConfirmIcon(), func() {
		w.Close()
	})
	cancelButton := widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
		w.Close()
	})
	nextButton := widget.NewButtonWithIcon("Next result", theme.NavigateNextIcon(), func() {
		if i+1 == result.TotalNumFound {
			statusText.SetText("This is the last result")
		} else {
			updateSearchResultContent(appCtxt, w, results, i+1, onConfirm)
		}
	})
	previousButton := widget.NewButtonWithIcon("Previous result", theme.NavigateBackIcon(), func() {
		if i == 0 {
			statusText.SetText("This is the first result")
		} else {
			updateSearchResultContent(appCtxt, w, results, i-1, onConfirm)
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
				confirmButton,
			),
		),
		nil,   // Left
		nil,   // Right
		image, // Center
	)

	// Set container to window
	w.SetContent(globalContainer)
}
