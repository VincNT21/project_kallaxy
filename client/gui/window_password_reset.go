package gui

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/VincNT21/kallaxy/client/context"
)

// This "Password Reset" approach is a "fake" one for local mode only
// But client and server endpoints could be used to send a real email to user in online mode, with slightly modifications

func showPasswordLostSecondaryWindow(appCtxt *context.AppContext) {

	secondaryWindow := fyne.CurrentApp().NewWindow("Password Lost")
	secondaryWindow.CenterOnScreen()
	secondaryWindow.Resize(fyne.NewSize(640, 480))

	passwordResetStep1(appCtxt, secondaryWindow)

	secondaryWindow.Show()
}

func passwordResetStep1(appCtxt *context.AppContext, parentWindow fyne.Window) {
	// First step: send email with link
	pageTitleText := canvas.NewText("Please enter your email", color.White)
	pageTitleText.Alignment = fyne.TextAlignCenter
	pageTitleText.TextSize = 16

	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("email")

	confirmButton := widget.NewButtonWithIcon("Send reset link", theme.MailComposeIcon(), func() {
		passwordReset, err := appCtxt.APIClient.Auth.SendPasswordResetLink(emailEntry.Text)
		if err != nil {
			dialog.ShowError(err, parentWindow)
		} else {
			passwordResetStep2(appCtxt, parentWindow, emailEntry.Text, passwordReset.Username, passwordReset.ResetToken)
		}

	})
	cancelButton := widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
		parentWindow.Close()
	})

	parentWindow.SetContent(container.NewBorder(
		pageTitleText,
		container.NewHBox(cancelButton, layout.NewSpacer(), confirmButton),
		nil, nil,
		container.NewVBox(layout.NewSpacer(), emailEntry, layout.NewSpacer()),
	))
}

func passwordResetStep2(appCtxt *context.AppContext, parentWindow fyne.Window, email, username, resetToken string) {
	// Second step: wait for email
	waitText := canvas.NewText(fmt.Sprintf("Email sent to %s ! Please wait to receive it", email), color.White)
	waitText.Alignment = fyne.TextAlignCenter
	waitText.TextSize = 16
	progress := widget.NewProgressBar()
	readButton := widget.NewButtonWithIcon("Read it !", theme.MailComposeIcon(), func() {
		passwordResetStep3(appCtxt, parentWindow, username, email, resetToken)
	})
	readButton.Disable()

	go func() {
		for i := 0.0; i < 1.0; i += 0.1 {
			time.Sleep(time.Millisecond * 500)
			progress.SetValue(i)
		}
		readButton.Enable()
	}()

	parentWindow.SetContent(container.NewBorder(
		waitText,
		container.NewHBox(layout.NewSpacer(), readButton, layout.NewSpacer()),
		nil, nil,
		container.NewVBox(layout.NewSpacer(), progress, layout.NewSpacer()),
	))
}

func passwordResetStep3(appCtxt *context.AppContext, parentWindow fyne.Window, username, email, resetToken string) {
	// Third step: read email
	splitted := strings.Split(email, "@")
	var domain string
	if len(splitted) == 2 {
		domain = splitted[1]
	} else {
		domain = "example.com"
	}

	fakeUrlText := canvas.NewText(fmt.Sprintf("https://www.%s", domain), color.White)
	fakeUrlText.Alignment = fyne.TextAlignCenter
	fakeUrlText.TextSize = 14
	line0 := container.NewBorder(
		customSeparatorForShelf(),
		customSeparatorForShelf(),
		customSeparatorForShelf(),
		customSeparatorForShelf(),
		fakeUrlText,
	)

	line1 := canvas.NewText("You've got a new email !", color.White)
	line1.TextSize = 14
	line1.TextStyle.Italic = true
	line1.TextStyle.Bold = true

	line2 := canvas.NewText("<< Dear //insert_username_here//,", color.White)
	line2.TextSize = 14

	line2b := canvas.NewText(fmt.Sprintf("It's a joke %s, I know it's you!", username), color.White)
	line2b.TextSize = 14

	line3 := canvas.NewText("First, let me say a warm thank for using my Kallaxy app.", color.White)
	line3.TextSize = 14

	line4 := canvas.NewText("I've heard you've lost your password.", color.White)
	line4.TextSize = 14

	line5 := canvas.NewText("How sad :'(", color.White)
	line5.TextSize = 14

	line6 := canvas.NewText("But you can create a new one", color.White)
	line6.TextSize = 14

	line7 := canvas.NewText("by clicking on the button below!", color.White)
	line7.TextSize = 14

	line8 := canvas.NewText("Best regards,", color.White)
	line8.TextSize = 14

	line9 := canvas.NewText("and see you soon in your Kallaxy ;)", color.White)
	line9.TextSize = 14

	line10 := canvas.NewText("VincNT21, for Project Kallaxy", color.White)
	line10.TextSize = 14
	line10.Alignment = fyne.TextAlignTrailing

	newPasswordButton := widget.NewButton("Create a new password", func() {
		passwordResetStep4(appCtxt, parentWindow, resetToken)
	})

	parentWindow.SetContent(container.NewBorder(line0, nil, nil, nil, container.NewVBox(
		layout.NewSpacer(),
		line1,
		layout.NewSpacer(),
		line2,
		line2b,
		layout.NewSpacer(),
		line3,
		layout.NewSpacer(),
		line4,
		line5,
		layout.NewSpacer(),
		line6,
		line7,
		layout.NewSpacer(),
		newPasswordButton,
		layout.NewSpacer(),
		line8,
		line9,
		layout.NewSpacer(),
		line10,
		layout.NewSpacer(),
	)))
}

func passwordResetStep4(appCtxt *context.AppContext, parentWindow fyne.Window, resetToken string) {
	pageTitleText := canvas.NewText("Please set a new password", color.White)
	pageTitleText.Alignment = fyne.TextAlignCenter
	pageTitleText.TextSize = 16

	newPasswordEntry := widget.NewPasswordEntry()
	confirmPasswordEntry := widget.NewPasswordEntry()

	passwordForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "New Password", Widget: newPasswordEntry},
			{Text: "Confirm password", Widget: confirmPasswordEntry},
		},
		// Forms buttons
		OnSubmit: func() {
			if newPasswordEntry.Text != confirmPasswordEntry.Text {
				dialog.ShowInformation("Error", "Both passwords are different", parentWindow)
			} else {
				_, err := appCtxt.APIClient.Auth.SetNewPassword(newPasswordEntry.Text, resetToken)
				if err != nil {
					dialog.ShowError(err, parentWindow)
				} else {
					dialog.ShowConfirm("Confirm", "New password set\nPlease login", func(b bool) {
						if b {
							appCtxt.PageManager.ShowLoginPage()
							parentWindow.Close()
						}
					}, parentWindow)
				}
			}
		},
		OnCancel: func() {
			parentWindow.Close()
		},
	}

	parentWindow.SetContent(container.NewBorder(
		pageTitleText,
		nil, nil, nil,
		container.NewVBox(layout.NewSpacer(), passwordForm, layout.NewSpacer()),
	))

}
