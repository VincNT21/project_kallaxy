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

func createEditRecordContent(appCtxt *context.AppContext, mediaType string, mediumWithRecord models.MediumWithRecord) *fyne.Container {
	// Create UI objects
	// Texts
	statusLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})
	pageTitleText := canvas.NewText(fmt.Sprintf("Update %s's record about %s (%s)", appCtxt.APIClient.CurrentUser.Username, mediumWithRecord.Title, mediaType), color.White)
	pageTitleText.TextSize = 20
	pageTitleText.Alignment = fyne.TextAlignCenter
	pageTitleText.TextStyle.Bold = true

	// User's Record forms
	// Date entries have a Action Button that calls a Date Picker dialog
	startDateEntry := widget.NewEntry()

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

	// Set initial date according to existing record
	if mediumWithRecord.StartDate == "" {
		startDateEntry.SetPlaceHolder("2025/01/01")
	} else {
		parsedDate, err := appCtxt.APIClient.Helpers.FormatDateToLocalFormat(mediumWithRecord.StartDate)
		if err != nil {
			startDateEntry.SetText("")
		} else {
			startDateEntry.SetText(parsedDate)
		}
	}

	endDateEntry := widget.NewEntry()

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

	// Set initial date according to existing record
	if mediumWithRecord.EndDate == "" {
		endDateEntry.SetPlaceHolder("2025/01/01")
	} else {
		parsedDate, err := appCtxt.APIClient.Helpers.FormatDateToLocalFormat(mediumWithRecord.EndDate)
		if err != nil {
			endDateEntry.SetText("")
		} else {
			endDateEntry.SetText(parsedDate)
		}
	}

	commentsEntry := widget.NewMultiLineEntry()
	commentsFormItem := widget.NewFormItem("Personal comments", commentsEntry)
	if mediumWithRecord.Comments != "" {
		commentsEntry.SetText(mediumWithRecord.Comments)
	}

	recordForm := widget.NewForm(startDateFormItem, endDateFormItem, commentsFormItem)

	// UI Buttons

	// Undo changes Button (to get back to existing data)
	buttonUndoChanges := widget.NewButtonWithIcon("Undo all changes", theme.ContentUndoIcon(), func() {
		dialog.ShowConfirm("Confirm", "Are you sure you want to undo all changes ?", func(b bool) {
			if b {
				parsedStartDate, err := appCtxt.APIClient.Helpers.FormatDateToLocalFormat(mediumWithRecord.StartDate)
				if err != nil {
					startDateEntry.SetText("")
				} else {
					startDateEntry.SetText(parsedStartDate)
				}

				parsedEndDate, err := appCtxt.APIClient.Helpers.FormatDateToLocalFormat(mediumWithRecord.EndDate)
				if err != nil {
					endDateEntry.SetText("")
				} else {
					endDateEntry.SetText(parsedEndDate)
				}

				commentsEntry.SetText(mediumWithRecord.Comments)
			}
		}, appCtxt.MainWindow)
	})

	exitButton := widget.NewButtonWithIcon("Homepage", theme.HomeIcon(), func() {
		dialog.ShowConfirm("Exit", "Are you sure you want to go back to Homepage ?\n\nAll unsubmitted changes will be lost!", func(b bool) {
			if b {
				appCtxt.PageManager.ShowHomePage()
			}
		}, appCtxt.MainWindow)
	})

	submitButton := widget.NewButtonWithIcon("Update", theme.ConfirmIcon(), func() {
		buttonFuncSubmitEditRecord(appCtxt, mediumWithRecord, startDateEntry, endDateEntry, commentsEntry)
	})

	// Group objects
	groupForms := container.NewVBox(widget.NewSeparator(), recordForm, widget.NewSeparator())
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

func buttonFuncSubmitEditRecord(appCtxt *context.AppContext, mediumWithRecord models.MediumWithRecord, startDateEntry, endDateEntry, commentsEntry *widget.Entry) {
	// Confirm info dialog box
	dialog.ShowCustomConfirm(
		"Confirm",
		"Create",
		"Cancel",
		container.NewVBox(
			widget.NewLabelWithStyle(fmt.Sprintf("Start Date: %s", startDateEntry.Text), fyne.TextAlignLeading, fyne.TextStyle{}),
			widget.NewLabelWithStyle(fmt.Sprintf("End Date: %s", endDateEntry.Text), fyne.TextAlignLeading, fyne.TextStyle{}),
			widget.NewLabelWithStyle(fmt.Sprintf("Comments: %s", commentsEntry.Text), fyne.TextAlignLeading, fyne.TextStyle{}),
		),
		func(b bool) {
			// If Confirmed. call the UpdateRecord client API function
			if b {
				_, err := appCtxt.APIClient.Records.UpdateRecord(
					mediumWithRecord.ID,
					startDateEntry.Text,
					endDateEntry.Text,
					commentsEntry.Text,
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
