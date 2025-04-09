package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
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

func customSpacerVertical(height float32) fyne.CanvasObject {
	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(1, height))
	return spacer
}

func createFallbackImage() fyne.CanvasObject {
	// Create a container with an icon and text
	brokenIcon := canvas.NewImageFromResource(theme.ErrorIcon())
	brokenIcon.SetMinSize(fyne.NewSize(50, 50))

	messageText := canvas.NewText("Image not available", color.White)
	messageText.Alignment = fyne.TextAlignCenter

	return container.NewCenter(
		container.NewVBox(
			brokenIcon,
			messageText,
		),
	)
}
