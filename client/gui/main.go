package gui

import (
	"fyne.io/fyne/v2/app"
)

func StartGUI() {
	a := app.New()
	getLoginWindow(a)

	a.Run()
}
