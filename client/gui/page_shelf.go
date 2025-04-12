package gui

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/VincNT21/kallaxy/client/context"
	"github.com/VincNT21/kallaxy/client/models"
)

func createShelfContent(appCtxt *context.AppContext, mediaRecords models.MediaWithRecords) *fyne.Container {
	// Create UI Objects
	// Texts
	pageTitleText := canvas.NewText(fmt.Sprintf("%s's Kallaxy Shelf", appCtxt.APIClient.CurrentUser.Username), color.White)
	pageTitleText.TextSize = 20
	pageTitleText.Alignment = fyne.TextAlignCenter
	pageTitleText.TextStyle.Bold = true

	// Buttons
	exitButton := widget.NewButtonWithIcon("Homepage", theme.HomeIcon(), func() {
		appCtxt.PageManager.ShowHomePage()
	})

	// Create the Shelf
	shelfContainer, err := buildMediaContainers(appCtxt, mediaRecords)
	if err != nil {
		// If an error occured while building the Shelf, return a valid emergency container
		errorContainer := container.NewVBox(
			widget.NewLabel("Error while constructing your shelf"),
			widget.NewButton("Return to home", func() {
				appCtxt.PageManager.ShowHomePage()
			}),
		)
		return container.NewBorder(
			container.NewVBox(pageTitleText),
			container.NewHBox(exitButton),
			nil, nil,
			errorContainer,
		)
	}

	// Create the global frame
	globalContainer := container.NewBorder(
		pageTitleText, // Top
		container.NewHBox(layout.NewSpacer(), exitButton), // Bottom
		customSpacerHorizontal(50),                        // Left
		customSpacerHorizontal(50),                        // Right
		shelfContainer,
	)

	return globalContainer
}

func buildMediaContainers(appCtxt *context.AppContext, mediaRecords models.MediaWithRecords) (*container.Scroll, error) {
	// This function create a scrollable shelf container where each media type has a compartment
	shelf := container.NewVBox()

	// Get media types map
	typesMap := appCtxt.APIClient.Helpers.GetMediaTypes(mediaRecords)

	// Iterate over each media type
	for mediaType := range typesMap {
		// Create the top separator
		topTextButton := widget.NewButton(strings.ToTitle(mediaType), func() {
			appCtxt.PageManager.ShowCompartmentTreePage(mediaType, mediaRecords.MediaRecords[mediaType])
		})
		topText := canvas.NewText(strings.ToTitle(mediaType), color.White)
		topText.Alignment = fyne.TextAlignCenter
		topText.TextSize = 20

		topSeparator := container.NewBorder(
			customSeparatorForShelf(),
			customSeparatorForShelf(),
			customSeparatorForShelf(),
			customSeparatorForShelf(),
			topTextButton,
		)

		// Create all images for media of this type into a Grid Wrap
		mediaDisplay := container.NewGridWrap(fyne.NewSize(200, 500))
		for _, medium := range mediaRecords.MediaRecords[mediaType] {
			// Get image buffer
			buffer, err := appCtxt.APIClient.Helpers.GetImage(medium.ImageUrl)
			if err != nil {
				return container.NewVScroll(shelf), err
			}

			image := canvas.NewImageFromReader(buffer, medium.Title)
			image.FillMode = canvas.ImageFillContain

			mediaDisplay.Add(image)
		}

		// Put the Grid Wrap inside a Border Container
		mediaCompartment := container.NewBorder(
			customSeparatorForShelf(),
			customSeparatorForShelf(),
			customSeparatorForShelf(),
			customSeparatorForShelf(),
			mediaDisplay,
		)

		// Add them to main shelf container
		shelf.Add(topSeparator)
		shelf.Add(mediaCompartment)
	}

	// Make the shelf scrollable
	scrollableShelf := container.NewVScroll(shelf)
	scrollableShelf.SetMinSize(fyne.NewSize(800, 600))

	// Return the completed shelf
	return scrollableShelf, nil
}

func createMediaTreeContent(appCtxt *context.AppContext, mediaType string, mediaList []models.MediumWithRecord) *fyne.Container {
	// Get the tree populated from helper function
	tree := createAndPopulateTree(appCtxt, mediaType, mediaList)

	// Make the tree scrollable
	scrollableTree := container.NewVScroll(tree)
	scrollableTree.SetMinSize(fyne.NewSize(800, 600))

	// Interface
	pageTitle := canvas.NewText(fmt.Sprintf("My Shelf Compartment - %s", mediaType), color.White)
	pageTitle.Alignment = fyne.TextAlignCenter
	pageTitle.TextSize = 18

	backButton := widget.NewButtonWithIcon("Back to Shelf", theme.ContentUndoIcon(), func() {
		appCtxt.PageManager.ShowShelfPage()
	})

	// Organize global container
	globalContainer := container.NewBorder(
		pageTitle,
		backButton,
		nil, nil,
		scrollableTree,
	)

	return globalContainer
}
