package gui

import "fyne.io/fyne/v2"

func (pm *GuiPageManager) GetHomeWindow() {
	// Create the window
	w := pm.appGui.NewWindow("Kallaxy")
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(800, 600))

	// Create objects

	// Set container to window
	w.Show()
}
