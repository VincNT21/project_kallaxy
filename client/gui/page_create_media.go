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
	"github.com/VincNT21/kallaxy/client/models"
	datepicker "github.com/sdassow/fyne-datepicker"
)

func (pm *GuiPageManager) GetCreateMediaWindow() {
	// Create the window
	w := pm.appGui.NewWindow("Kallaxy")
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(800, 600))

	// Create texts
	statusLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})
	pageTitleText := canvas.NewText(fmt.Sprintf("Create a new %s", pm.mediaType), color.White)
	pageTitleText.TextSize = 20
	pageTitleText.Alignment = fyne.TextAlignCenter
	pageTitleText.TextStyle.Bold = true

	// Create objects
	titleEntry := widget.NewEntry()
	titleForm := widget.NewFormItem("Title", titleEntry)

	mediaTypeEntry := widget.NewEntry()
	mediaTypeEntry.SetText(pm.mediaType)
	mediaTypeEntry.Disable()
	mediaTypeForm := widget.NewFormItem("Media Type", mediaTypeEntry)

	creatorEntry := widget.NewEntry()
	creatorForm := widget.NewFormItem("Creator", creatorEntry)

	releaseYearEntry := widget.NewEntry()
	releaseYearForm := widget.NewFormItem("Release Year", releaseYearEntry)

	imageUrlEntry := widget.NewEntry()
	imageUrlForm := widget.NewFormItem("Image URL", imageUrlEntry)

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
			w,
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
			w,
		)
	})
	endDateForm := widget.NewFormItem("Completed on", endDateEntry)

	exitButton := widget.NewButtonWithIcon("Homepage", theme.HomeIcon(), func() {
		dialog.ShowConfirm("Exit", "Are you sure you want to go back to Homepage ?\n\nAll unsubmitted changes will be lost!", func(b bool) {
			if b {
				pm.GetHomeWindow()
				w.Close()
			}
		}, w)
	})

	submitButton := widget.NewButtonWithIcon("Submit", theme.ConfirmIcon(), func() {
		// Confirm info dialog box
		dialog.ShowCustomConfirm(
			"Confirm",
			"Create",
			"Cancel",
			container.NewVBox(
				widget.NewLabelWithStyle("Do you confirm the info entered ?", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
				widget.NewLabelWithStyle(fmt.Sprintf("Media Type: %s", mediaTypeEntry.Text), fyne.TextAlignCenter, fyne.TextStyle{}),
				widget.NewLabelWithStyle(fmt.Sprintf("Title: %s", titleEntry.Text), fyne.TextAlignCenter, fyne.TextStyle{}),
				widget.NewLabelWithStyle(fmt.Sprintf("Creator: %s", creatorEntry.Text), fyne.TextAlignCenter, fyne.TextStyle{}),
				widget.NewLabelWithStyle(fmt.Sprintf("Release Year: %s", releaseYearEntry.Text), fyne.TextAlignCenter, fyne.TextStyle{}),
				widget.NewLabelWithStyle(fmt.Sprintf("Image URL: %s", imageUrlEntry.Text), fyne.TextAlignCenter, fyne.TextStyle{}),
				widget.NewLabelWithStyle(fmt.Sprintf("Start Date: %s", startDateEntry.Text), fyne.TextAlignCenter, fyne.TextStyle{}),
				widget.NewLabelWithStyle(fmt.Sprintf("End Date: %s", endDateEntry.Text), fyne.TextAlignCenter, fyne.TextStyle{}),
			),
			func(b bool) {
				// If Confirmed
				if b {
					_, _, err := pm.appCtxt.APIClient.Media.CreateMediumAndRecord(
						titleEntry.Text,
						mediaTypeEntry.Text,
						creatorEntry.Text,
						releaseYearEntry.Text,
						imageUrlEntry.Text,
						startDateEntry.Text,
						endDateEntry.Text,
					)
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
							dialog.ShowInformation("Error", "There is a problem with your request:\n- One field is missing in the form\nAND/OR\n- Start date is before end date", w)
						case models.ErrConflict:
							dialog.ShowInformation("Error", "A medium with the same couple title & media type already exists", w)
						case models.ErrNotFound:
							dialog.ShowInformation("Error", "Not found", w)
						default:
							dialog.ShowError(err, w)
						}
					} else {
						log.Println("--GUI-- CreateMedia Form successful")
						pm.GetHomeWindow()
						w.Close()
					}
				}
			}, w,
		)
	})

	// Group objects
	mediaForm := widget.NewForm(titleForm, mediaTypeForm, creatorForm, releaseYearForm, imageUrlForm)
	recordForm := widget.NewForm(startDateForm, endDateForm)
	groupForms := container.NewVBox(mediaForm, widget.NewSeparator(), recordForm)
	submitRow := container.NewHBox(layout.NewSpacer(), submitButton, layout.NewSpacer())
	statusRow := container.NewHBox(layout.NewSpacer(), statusLabel, layout.NewSpacer())
	bottomRow := container.NewHBox(layout.NewSpacer(), exitButton)
	centralPart := container.NewVBox(groupForms, statusRow, submitRow)

	// Create the global frame
	globalContainer := container.NewBorder(pageTitleText, bottomRow, nil, nil, centralPart)

	// Set container to window
	w.SetContent(globalContainer)
	w.Show()
}
