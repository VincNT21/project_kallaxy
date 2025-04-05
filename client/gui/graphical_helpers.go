package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

func customSeparatorForShelf() *canvas.Rectangle {
	myColor := color.RGBA{R: 128, G: 0, B: 128, A: 255}
	separator := canvas.NewRectangle(myColor)
	separator.SetMinSize(fyne.NewSize(2, 2))
	return separator
}

func customSpacerHorizontal(width float32) fyne.CanvasObject {
	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(width, 1))
	return spacer
}
