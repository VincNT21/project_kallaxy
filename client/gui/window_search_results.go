package gui

import (
	"errors"
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/VincNT21/kallaxy/client/context"
	"github.com/VincNT21/kallaxy/client/models"
)

func initSearchResultContent(appCtxt *context.AppContext, parentWindow fyne.Window, mediumTitle, mediumType, vgPlatform string, entryMap map[string]*widget.Entry, onConfirm func(models.ClientMedium)) {
	// Create the window
	secondaryWindow := fyne.CurrentApp().NewWindow("Search Results")
	secondaryWindow.CenterOnScreen()
	secondaryWindow.Resize(fyne.NewSize(800, 600))

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
	updateSearchResultContent(appCtxt, mediumType, secondaryWindow, results, 0, entryMap, onConfirm)

	// Display window
	secondaryWindow.Show()
}

func updateSearchResultContent(appCtxt *context.AppContext, mediaType string, w fyne.Window, results []models.ShortOnlineSearchResult, i int, entryMap map[string]*widget.Entry, onConfirm func(models.ClientMedium)) {

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

	// Inside function lo load image
	var imageObj fyne.CanvasObject
	loadImage := func() fyne.CanvasObject {
		// Fetch the image as a buffer
		imageUrl := result.ImageUrl
		if mediaType == "boardgame" {
			actualImageUrl, err := appCtxt.APIClient.Helpers.GetBoardgameImageUrl(result.ApiID)
			fmt.Println(actualImageUrl)
			if err != nil {
				statusText.SetText("Could not find image URL")
				return createFallbackImage()
			}
			imageUrl = actualImageUrl
			result.ImageUrl = actualImageUrl
		}

		bufImage, err := appCtxt.APIClient.Helpers.GetImage(imageUrl)
		if err != nil {
			statusText.SetText(fmt.Sprintf("Error loading image: %v\n", err))
			return createFallbackImage()
		}
		// Create the image component
		image := canvas.NewImageFromReader(bufImage, "image")
		image.FillMode = canvas.ImageFillContain
		image.SetMinSize(fyne.NewSize(350, 250))

		return image
	}

	imageObj = loadImage()

	// Buttons
	detailsButton := widget.NewButtonWithIcon("Get details", theme.SearchIcon(), func() {
		showSearchMediumDetails(appCtxt, mediaType, result.ApiID, result.ImageUrl, w, entryMap, onConfirm)

	})
	cancelButton := widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
		w.Close()
	})
	nextButton := widget.NewButtonWithIcon("Next result", theme.NavigateNextIcon(), func() {
		if i+1 == result.TotalNumFound {
			statusText.SetText("This is the last result")
		} else {
			// Show next page of results
			updateSearchResultContent(appCtxt, mediaType, w, results, i+1, entryMap, onConfirm)
		}
	})
	previousButton := widget.NewButtonWithIcon("Previous result", theme.NavigateBackIcon(), func() {
		if i == 0 {
			statusText.SetText("This is the first result")
		} else {
			// Show previous page of results
			updateSearchResultContent(appCtxt, mediaType, w, results, i-1, entryMap, onConfirm)
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
		nil,      // Left
		nil,      // Right
		imageObj, // Center
	)

	// Set container to window
	w.SetContent(globalContainer)
}

func showSearchMediumDetails(appCtxt *context.AppContext, mediaType, mediumApiID, imageUrl string, parentWindow fyne.Window, entryMap map[string]*widget.Entry, onConfirm func(models.ClientMedium)) {
	// Get details for medium on external API
	medium, err := appCtxt.APIClient.Helpers.SearchMediumDetailsOnExternalApi(mediaType, mediumApiID)
	if err != nil {
		dialog.ShowError(fmt.Errorf("couldn't get details about medium: %v", err), parentWindow)
		return
	}

	var urlImage string
	if mediaType == "book" && parentWindow == appCtxt.MainWindow {
		// This means that showSearchMediumDetails is called for a "Search by ISBN"
		// No imageUrl is provided, need to find it online
		urlImage = fmt.Sprintf("https://covers.openlibrary.org/b/isbn/%s-M.jpg", mediumApiID)
	} else {
		urlImage = imageUrl
	}

	// Add image Url and set empty fields to "unknown"
	medium.ImageUrl = urlImage // Insert back the proper image url
	if medium.Creator == "" {
		medium.Creator = "unknown"
	}
	if medium.PubDate == "" {
		medium.PubDate = "unknown"
	}

	// Prepare results
	titleText := canvas.NewText(fmt.Sprintf("Title: %s", medium.Title), color.White)
	titleText.TextSize = 16
	creatorText := canvas.NewText(fmt.Sprintf("Creator: %s", medium.Creator), color.White)
	creatorText.TextSize = 16
	pubDateText := canvas.NewText(fmt.Sprintf("Publication date: %s", medium.PubDate), color.White)
	pubDateText.TextSize = 16

	metadataBox := createMetadataTextContainer(appCtxt, entryMap, medium)

	// Inside function lo load image
	var imageObj fyne.CanvasObject
	loadImage := func() fyne.CanvasObject {
		// Fetch the image as a buffer
		bufImage, err := appCtxt.APIClient.Helpers.GetImage(urlImage)
		if err != nil {
			return createFallbackImage()
		}
		// Create the image component
		image := canvas.NewImageFromReader(bufImage, "image")
		image.FillMode = canvas.ImageFillContain
		image.SetMinSize(fyne.NewSize(350, 250))

		return image
	}

	imageObj = loadImage()

	// Display them in a dialog box
	dialog.ShowCustomConfirm(
		"Details",
		"Confirm",
		"Dismiss",
		container.NewVBox(imageObj, titleText, creatorText, pubDateText, metadataBox),
		func(b bool) {
			if b {
				// If user confirms, call OnConfirm callback function
				onConfirm(medium)

				if parentWindow != appCtxt.MainWindow {
					parentWindow.Close()
				}
			}
		},
		parentWindow,
	)

}
